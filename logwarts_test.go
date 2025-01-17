package logwarts_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"math/rand"
	"sync"
	"testing"
	"time"

	"lighthouse/internal/services/logwarts"
)

func testLogger(t *testing.T, logger *logwarts.Logger, buf *bytes.Buffer, defaultAttributes map[string]any) {
	// Create a context and populate it with defaultAttributes
	ctx := context.Background()
	for key, value := range defaultAttributes {
		ctx = context.WithValue(ctx, key, value)
	}

	logger.Info(ctx, "Info message", slog.String("key", "value"))
	logger.Warn(ctx, "Warn message")
	logger.Debug(ctx, "Debug message")
	logger.Error(ctx, "Error message", slog.Int("errorCode", 500))

	var logEntries []map[string]any
	decoder := json.NewDecoder(buf)
	for {
		var logEntry map[string]any
		if err := decoder.Decode(&logEntry); err != nil {
			if err.Error() != "EOF" {
				t.Fatalf("Failed to decode log entry: %v", err)
			}
			break
		}
		logEntries = append(logEntries, logEntry)
	}

	tests := []struct {
		level    string
		message  string
		extraKey string
		extraVal any
	}{
		{"INFO", "Info message", "key", "value"},
		{"WARN", "Warn message", "", nil},
		{"DEBUG", "Debug message", "", nil},
		{"ERROR", "Error message", "errorCode", float64(500)},
	}

	for _, test := range tests {
		found := false
		for _, entry := range logEntries {
			if entry["level"] == test.level && entry["msg"] == test.message {
				found = true
				// Check that all defaultAttributes are present in the log entry
				for key, expectedValue := range defaultAttributes {
					if entry[key] != expectedValue {
						t.Errorf("Expected %s=%v in log, got %v. Full log: %+v", key, expectedValue, entry[key], entry)
					}
				}
				if test.extraKey != "" && entry[test.extraKey] != test.extraVal {
					t.Errorf("Expected %s=%v in log for '%s', got %v. Full log: %+v", test.extraKey, test.extraVal, test.message, entry[test.extraKey], entry)
				}
			}
		}
		if !found {
			t.Errorf("Log entry for '%s' with level '%s' not found. Logs: %+v", test.message, test.level, logEntries)
		}
	}
}

func TestGetLogger(t *testing.T) {
	defaultKeys := []string{
		"userID",
		"requestID",
		"name",
		"email",
		"phone",
		"address",
		"city",
		"state",
		"zip",
		"country",
		"age",
		"gender",
		"accountType",
		"subscriptionStatus",
		"signupDate",
		"lastLogin",
		"ipAddress",
		"deviceType",
		"browser",
		"os",
		"language",
		"timezone",
		"preferredCurrency",
		"purchaseCount",
		"totalSpent",
		"membershipLevel",
		"referralCode",
		"loyaltyPoints",
		"notificationsEnabled",
		"emailVerified",
		"phoneVerified",
		"twoFactorAuth",
		"marketingOptIn",
		"birthdate",
		"favoriteCategory",
		"wishlistItems",
		"cartItems",
		"savedCards",
		"defaultPaymentMethod",
		"shippingAddress",
		"billingAddress",
		"preferredShippingMethod",
		"orderHistory",
		"averageOrderValue",
		"accountNotes",
		"customerSupportTickets",
		"supportRating",
		"favoriteBrand",
		"promoCodesUsed",
		"returnRate",
	}

	buf := &bytes.Buffer{}
	logger := slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelDebug, // Include DEBUG level logs in the output
	}))

	customLogger := logwarts.GetLogger(defaultKeys, logger)

	defaultAttributes := map[string]any{
		"userID":                  "42069",
		"requestID":               "abcd-efgh",
		"name":                    "Big Chungus",
		"email":                   "chungus@looneytunes.com",
		"phone":                   "+4204204200",
		"address":                 "123 Meme Blvd, Danktown",
		"city":                    "Memetropolis",
		"state":                   "LOL",
		"zip":                     "1337",
		"country":                 "Neverland",
		"age":                     "69",
		"gender":                  "MemeLord",
		"accountType":             "Ultra Premium Deluxe Platinum Plus",
		"subscriptionStatus":      "Subscribed to Chaos",
		"signupDate":              "1969-04-20",
		"lastLogin":               "Never Logs Out",
		"ipAddress":               "127.0.0.1",
		"deviceType":              "Toaster",
		"browser":                 "Internet Explorer",
		"os":                      "Windows 95",
		"language":                "Emoji",
		"timezone":                "UTC+420",
		"preferredCurrency":       "Dogecoin",
		"purchaseCount":           "420",
		"totalSpent":              "69.69",
		"membershipLevel":         "Supreme Overlord",
		"referralCode":            "CHUNGUS4LIFE",
		"loyaltyPoints":           "9001",
		"notificationsEnabled":    "false",
		"emailVerified":           "false",
		"phoneVerified":           "false",
		"twoFactorAuth":           "LOL, no",
		"marketingOptIn":          "true",
		"birthdate":               "2000-01-01",
		"favoriteCategory":        "Unnecessary Gadgets",
		"wishlistItems":           "69",
		"cartItems":               "4",
		"savedCards":              "0",
		"defaultPaymentMethod":    "Magic Beans",
		"shippingAddress":         "The Moon",
		"billingAddress":          "Under the Couch",
		"preferredShippingMethod": "Pigeon Express",
		"orderHistory":            "9999",
		"averageOrderValue":       "420.42",
		"accountNotes":            "Frequently buys weird stuff at 3 AM",
		"customerSupportTickets":  "42",
		"supportRating":           "2.5,",
		"favoriteBrand":           "YeetCo",
		"promoCodesUsed":          "666",
		"returnRate":              "0.01,",
	}

	testLogger(t, &customLogger, buf, defaultAttributes)

}

func TestHighFrequencyLogging(t *testing.T) {

	numberOfLogs := 9999
	// Setup default context keys
	defaultKeys := []string{"userID", "requestID"}

	// Redirect logger output to a buffer for testing
	buf := &bytes.Buffer{}
	logger := slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	customLogger := logwarts.GetLogger(defaultKeys, logger)

	// Create a context with values for the default keys
	ctx := context.WithValue(context.Background(), "userID", "12345")
	ctx = context.WithValue(ctx, "requestID", "abcd-efgh")

	// Define log levels and messages
	logLevels := []slog.Level{slog.LevelInfo, slog.LevelWarn, slog.LevelDebug, slog.LevelError}
	logMessages := []string{"Info log", "Warn log", "Debug log", "Error log"}

	// Create a new random source and generator
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate 500 log calls in quick succession
	var wg sync.WaitGroup
	wg.Add(numberOfLogs)

	for i := 0; i < numberOfLogs; i++ {
		go func() {
			defer wg.Done()
			level := logLevels[rnd.Intn(len(logLevels))]
			message := logMessages[rnd.Intn(len(logMessages))]
			attributes := []slog.Attr{slog.Int("randomID", rnd.Intn(1000))}

			switch level {
			case slog.LevelInfo:
				customLogger.Info(ctx, message, attributes...)
			case slog.LevelWarn:
				customLogger.Warn(ctx, message, attributes...)
			case slog.LevelDebug:
				customLogger.Debug(ctx, message, attributes...)
			case slog.LevelError:
				customLogger.Error(ctx, message, attributes...)
			}
		}()
	}

	// Wait for all log calls to complete
	wg.Wait()

	// Parse the logs from the buffer
	var logEntries []map[string]any
	decoder := json.NewDecoder(buf)
	for {
		var logEntry map[string]any
		if err := decoder.Decode(&logEntry); err != nil {
			if err.Error() != "EOF" {
				t.Fatalf("Failed to decode log entry: %v", err)
			}
			break
		}
		logEntries = append(logEntries, logEntry)
	}

	// Verify the number of log entries matches the expected count
	if len(logEntries) != numberOfLogs {
		t.Errorf("Expected %d log entries, got %d", numberOfLogs, len(logEntries))
	}

	// Validate that each entry contains required fields
	for _, entry := range logEntries {
		if entry["userID"] != "12345" || entry["requestID"] != "abcd-efgh" {
			t.Errorf("Context values missing or incorrect in log: %+v", entry)
		}
		if _, ok := entry["randomID"]; !ok {
			t.Errorf("Random ID missing in log: %+v", entry)
		}
	}
}
