package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gosuri/uilive"
	"github.com/spf13/cobra"

	"github.com/parthrs/btctracker/pkg/coinbaseutils"
	"github.com/parthrs/btctracker/pkg/common"
)

// Color name to code mapping.
var colorMap = map[string]string{
	"red":   "\033[31m",
	"green": "\033[32m",
	"reset": "\033[0m",
}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "btctracker",
	Short: "Track BTC price in real-time.",
	Long:  `This command tracks the value of BTC in real-time right from your terminal.`,

	Run: func(cmd *cobra.Command, args []string) {
		interval, _ := cmd.Flags().GetInt("interval")
		bucket, _ := cmd.Flags().GetInt("bucket")

		// Channel to listen for user interrupts (Ctrl+c)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		trackBtc(interval, bucket, c)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.btctracker.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().Int("interval", 10, "Interval between each refresh")
	rootCmd.Flags().Int("bucket", 10, "Size of the bucket (number of price data points) to print")
}

// Connecting the BTC-USD price, time stamp at the time of price (approx)
// and the color corresponding to an increase or decrease (green and red respectively)
// into a single type.
type BtcPoint struct {
	Price float64
	Time  string
	Color string
}

// NewBtcPoint returns type BtcPoint, it takes the price, timestamp and
// the color (as string) as input.
func NewBtcPoint(price float64, time, color string) BtcPoint {
	return BtcPoint{
		Price: price,
		Time:  time,
		Color: color,
	}
}

// printScreen paints the stdout in real-time by iterating over the linked list
// and printing the Value of each node (type BtcPoint). It prints horizontally
// by first printing the price and then the timestamp corresponding to the price.
func printScreen(l *common.SinglyLinkedList[BtcPoint], w *uilive.Writer) {

	// Print the first row which is the prices for each point
	fmt.Fprintf(w, "%-7s", "Price")

	// Initialize the head node before looping
	node := l.Head

	for i := 0; i < l.Capacity; i++ {
		// Price of zero (0.00) indicates an un-initialized BtcPoint node.
		// Print a "-" and continue.
		// Similarly a price of -1 indicates error in fetching from the Coinbase api.
		if node.Value.Price == float64(0) {
			fmt.Fprintf(w, "%-10s", "-")
			node = node.Next
			continue
		} else if node.Value.Price == float64(-1) {
			fmt.Fprintf(w, "%-10s", "Err!")
			node = node.Next
			continue
		}

		// Set the color before price is printed.
		if node.Value.Color != "" {
			fmt.Fprintf(w, colorMap[node.Value.Color])
		}

		// Print with a precision of 2 decimal points and a padding
		// of 10 spaces.
		fmt.Fprintf(w, "%-10.2f", node.Value.Price)
		fmt.Fprintf(w, colorMap["reset"])

		node = node.Next
	}

	fmt.Fprintf(w, "\n")

	// Print the second row which is the timestamps for the prices.
	fmt.Fprintf(w, "%-7s", "Time")
	node = l.Head
	for i := 0; i < l.Capacity; i++ {
		fmt.Fprintf(w, "%-10s", node.Value.Time)
		node = node.Next
	}
	fmt.Fprintf(w, "\n")
}

// trackBtc maintains a fixed size queue (which is set by the user via the bucket param)
// denoting the number of price, timestamp columns or BTC-USD points. Each new price fetch
// is added to the queue (removing the oldest element to get removed) and printScreen is called
// post.
func trackBtc(interval, bucket int, interrupt chan os.Signal) {

	// Setup the stdout printer.
	printer := uilive.New()
	printer.Start()
	defer printer.Stop()

	// Initialize the list with "placeholder" values.
	l := common.NewSinglyLinkedList[BtcPoint](nil, nil, bucket)
	for i := 0; i < bucket-1; i++ {
		l.AddNodeAtTail(NewBtcPoint(float64(0), "-", ""))
	}

	// Get initial BTC-USD value and print to stdout.
	price, err := coinbaseutils.GetBtcUsdPrice()
	if err != nil {
		price = -1
	}
	h, m, s := time.Now().Clock()
	l.AddNodeAtTail(NewBtcPoint(price, fmt.Sprintf("%02d:%02d:%02d", h, m, s), ""))
	printScreen(l, printer)

	// Creates a trigger every interval seconds to poll the latest BTC-USD
	// value and refresh stdout.
	ticker := time.NewTicker(time.Second * time.Duration(interval))
	defer ticker.Stop()

	// Continous loop to either listen for a user interrupt or the timer trigger
	for {
		select {
		case <-interrupt:
			// Graceful exit; defer statements run to release resources
			return
		case <-ticker.C:
			newPrice, err := coinbaseutils.GetBtcUsdPrice()
			if err != nil {
				newPrice = -1
			}

			h, m, s := time.Now().Clock()

			if newPrice > price {
				l.AddNodeAtTail(NewBtcPoint(newPrice, fmt.Sprintf("%02d:%02d:%02d", h, m, s), "green"))
			} else if newPrice < price {
				l.AddNodeAtTail(NewBtcPoint(newPrice, fmt.Sprintf("%02d:%02d:%02d", h, m, s), "red"))
			} else {
				l.AddNodeAtTail(NewBtcPoint(newPrice, fmt.Sprintf("%02d:%02d:%02d", h, m, s), ""))
			}

			price = newPrice
			printScreen(l, printer)
		default:
			continue
		}
	}
}
