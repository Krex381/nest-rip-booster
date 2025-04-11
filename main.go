package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/inancgumus/screen"
	"github.com/muesli/termenv"
)

// KREX NEST.RIP BOOSTER - Copyright © 2025

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.4 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:99.0) Gecko/20100101 Firefox/99.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 15_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.4 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 15_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.4 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36 Edg/100.0.1185.50",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36 OPR/86.0.4363.64",
	"Mozilla/5.0 (Android 12; Mobile; rv:99.0) Gecko/99.0 Firefox/99.0",
	"Mozilla/5.0 (Android 12; Mobile; LG-M255; rv:99.0) Gecko/99.0 Firefox/99.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:100.0) Gecko/20100101 Firefox/100.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 12_3_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.4 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36 Edg/101.0.1210.32",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 15_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/101.0.4951.44 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Linux; Android 12) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.41 Mobile Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:100.0) Gecko/20100101 Firefox/100.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.59",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36 Edg/92.0.902.73",
}

var referers = []string{
	"https://www.google.com/",
	"https://www.bing.com/",
	"https://www.facebook.com/",
	"https://twitter.com/",
	"https://www.instagram.com/",
	"https://www.reddit.com/",
	"https://www.youtube.com/",
	"https://www.pinterest.com/",
	"https://www.linkedin.com/",
	"https://discord.com/",
	"https://duckduckgo.com/",
	"https://www.yahoo.com/",
	"https://www.twitch.tv/",
	"https://github.com/",
	"https://stackoverflow.com/",
}

var (
	successColor   = color.New(color.FgGreen, color.Bold)
	errorColor     = color.New(color.FgRed, color.Bold)
	highlightColor = color.New(color.FgCyan, color.Bold)
	warningColor   = color.New(color.FgYellow, color.Bold)
	infoColor      = color.New(color.FgWhite)
	titleColor     = color.New(color.FgHiMagenta, color.Bold)
	subtitleColor  = color.New(color.FgHiCyan)
	numberColor    = color.New(color.FgHiGreen, color.Bold)
	urlColor       = color.New(color.FgHiBlue, color.Underline)
	p              = termenv.ColorProfile()
)

const (
	logo = `
    ███▄    █ ▓█████   ██████ ▄▄▄█████▓            ██▀███   ██▓ ██▓███      ▄▄▄▄    ▒█████   ▒█████    ██████ ▄▄▄█████▓▓█████  ██▀███  
    ██ ▀█   █ ▓█   ▀ ▒██    ▒ ▓  ██▒ ▓▒           ▓██ ▒ ██▒▓██▒▓██░  ██▒   ▓█████▄ ▒██▒  ██▒▒██▒  ██▒▒██    ▒ ▓  ██▒ ▓▒▓█   ▀ ▓██ ▒ ██▒
   ▓██  ▀█ ██▒▒███   ░ ▓██▄   ▒ ▓██░ ▒░           ▓██ ░▄█ ▒▒██▒▓██░ ██▓▒   ▒██▒ ▄██▒██░  ██▒▒██░  ██▒░ ▓██▄   ▒ ▓██░ ▒░▒███   ▓██ ░▄█ ▒
   ▓██▒  ▐▌██▒▒▓█  ▄   ▒   ██▒░ ▓██▓ ░            ▒██▀▀█▄  ░██░▒██▄█▓▒ ▒   ▒██░█▀  ▒██   ██░▒██   ██░  ▒   ██▒░ ▓██▓ ░ ▒▓█  ▄ ▒██▀▀█▄  
   ▒██░   ▓██░░▒████▒▒██████▒▒  ▒██▒ ░     ██▓    ░██▓ ▒██▒░██░▒██▒ ░  ░   ░▓█  ▀█▓░ ████▓▒░░ ████▓▒░▒██████▒▒  ▒██▒ ░ ░▒████▒░██▓ ▒██▒
   ░ ▒░   ▒ ▒ ░░ ▒░ ░▒ ▒▓▒ ▒ ░  ▒ ░░       ▒▓▒    ░ ▒▓ ░▒▓░░▓  ▒▓▒░ ░  ░   ░▒▓███▀▒░ ▒░▒░▒░ ░ ▒░▒░▒░ ▒ ▒▓▒ ▒ ░  ▒ ░░   ░░ ▒░ ░░ ▒▓ ░▒▓░
   ░ ░░   ░ ▒░ ░ ░  ░░ ░▒  ░ ░    ░        ░▒       ░▒ ░ ▒░ ▒ ░░▒ ░        ▒░▒   ░   ░ ▒ ▒░   ░ ▒ ▒░ ░ ░▒  ░ ░    ░     ░ ░  ░  ░▒ ░ ▒░
      ░   ░ ░    ░   ░  ░  ░    ░          ░        ░░   ░  ▒ ░░░           ░    ░ ░ ░ ░ ▒  ░ ░ ░ ▒  ░  ░  ░    ░         ░     ░░   ░ 
            ░    ░  ░      ░                ░        ░      ░               ░          ░ ░      ░ ░        ░              ░  ░   ░     
                                            ░                                    ░                                                     
`
	copyright = `Copyright © 2025 KREX. All rights reserved.`

	// Window dimensions
	windowWidth  = 1260
	windowHeight = 637
)

type RequestResult struct {
	Success bool
	Error   error
}

func getRandomIP() string {
	return fmt.Sprintf("%d.%d.%d.%d",
		rand.Intn(223)+1,
		rand.Intn(255),
		rand.Intn(255),
		rand.Intn(255))
}

func getRandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

func getRandomReferer() string {
	return referers[rand.Intn(len(referers))]
}

// Set window title and size
func setWindowTitleAndSize(title string) {
	if runtime.GOOS == "windows" {
		// First set the window size
		// Convert pixel dimensions to columns and rows (approximate conversion)
		// Standard terminal font is about 8x16 pixels
		cols := windowWidth / 8
		rows := windowHeight / 16

		sizeCmd := exec.Command("cmd", "/c", "mode", "con:", fmt.Sprintf("cols=%d", cols), fmt.Sprintf("lines=%d", rows))
		_ = sizeCmd.Run()

		// Then set the title
		titleCmd := exec.Command("cmd", "/c", "title", title)
		_ = titleCmd.Run()
	} else {
		// Set title for Unix systems
		fmt.Printf("\033]0;%s\007", title)

		// For Unix systems, try to set terminal size using ANSI escape sequences
		// This is an approximation and may not work on all terminals
		fmt.Printf("\033[8;%d;%dt", windowHeight/16, windowWidth/8)
	}
}

func displayProgressBar(current, total int32, width int) string {
	percent := float64(current) / float64(total)
	done := int(percent * float64(width))

	bar := "["
	for i := 0; i < width; i++ {
		if i < done {
			bar += "█" // Full block character
		} else {
			bar += "░" // Light shade character
		}
	}
	bar += "]"

	return bar
}

func clickWorker(url string, results chan<- RequestResult, isUniqueClick bool, workerID int, uniqueClicks int, stopAt *int32) {
	transport := &http.Transport{
		Proxy: nil,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		MaxConnsPerHost:       100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableCompression:    true,
		DisableKeepAlives:     false,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	if atomic.LoadInt32(stopAt) <= 0 {
		results <- RequestResult{Success: false, Error: fmt.Errorf("target reached")}
		return
	}

	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		// Recheck target count
		if atomic.LoadInt32(stopAt) <= 0 {
			results <- RequestResult{Success: false, Error: fmt.Errorf("target reached")}
			return
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			time.Sleep(time.Duration(50+rand.Intn(200)) * time.Millisecond)
			continue
		}

		req.Header.Set("User-Agent", getRandomUserAgent())
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Referer", getRandomReferer())
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Upgrade-Insecure-Requests", "1")
		req.Header.Set("Cache-Control", "max-age=0")

		if isUniqueClick {
			randomSeed := 0
			if uniqueClicks > 0 && workerID < uniqueClicks {
				randomSeed = workerID
			} else {
				randomSeed = rand.Intn(10000)
			}

			ipBase := randomSeed % 100
			req.Header.Set("X-Forwarded-For", fmt.Sprintf("192.168.%d.%d", ipBase/10, ipBase%10+100))
			req.Header.Set("X-Real-IP", getRandomIP())

			if randomSeed%2 == 0 {
				req.Header.Set("DNT", "1")
			} else {
				req.Header.Set("DNT", "0")
			}

			width := 800 + (randomSeed % 1200)
			height := 600 + (randomSeed % 900)
			colorDepth := []string{"24", "32", "48"}[randomSeed%3]
			platform := []string{"Windows", "Linux", "macOS"}[randomSeed%3]
			browser := []string{"Chrome", "Edge", "Firefox"}[randomSeed%3]

			req.Header.Set("Sec-CH-UA-Platform", platform)
			req.Header.Set("Sec-CH-UA-Platform-Version", fmt.Sprintf("%d.%d.%d", 10+(randomSeed%3), randomSeed%5, randomSeed%100))
			req.Header.Set("Sec-CH-UA", fmt.Sprintf("\"%s\";v=\"%d\"", browser, 90+(randomSeed%11)))
			req.Header.Set("Viewport-Width", fmt.Sprintf("%d", width))
			req.Header.Set("Device-Memory", []string{"4", "8", "16"}[randomSeed%3])
			req.Header.Add("Cookie", fmt.Sprintf("uid=%d; resolution=%dx%d; colorDepth=%s", randomSeed, width, height, colorDepth))
		}

		resp, err := client.Do(req)
		if err != nil {
			time.Sleep(time.Duration(100+rand.Intn(500)) * time.Millisecond)
			continue
		}

		_, _ = io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			curCount := atomic.LoadInt32(stopAt)
			if curCount > 0 {
				if atomic.CompareAndSwapInt32(stopAt, curCount, curCount-1) {
					results <- RequestResult{Success: true, Error: nil}
					return
				}
			}
			results <- RequestResult{Success: false, Error: fmt.Errorf("target reached")}
			return
		}

		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			time.Sleep(time.Duration(500+rand.Intn(1000)) * time.Millisecond)
		} else {
			time.Sleep(time.Duration(100+rand.Intn(300)) * time.Millisecond)
		}
	}

	results <- RequestResult{Success: false, Error: fmt.Errorf("all retries failed")}
}

func animateLogo() {
	colors := []termenv.Color{
		p.Color("#ff0000"), // Red
		p.Color("#ff7f00"), // Orange
		p.Color("#ffff00"), // Yellow
		p.Color("#00ff00"), // Green
		p.Color("#0000ff"), // Blue
		p.Color("#4b0082"), // Indigo
		p.Color("#9400d3"), // Violet
	}

	lines := strings.Split(logo, "\n")
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()

	for i := 0; i < 10; i++ {
		screen.Clear()
		screen.MoveTopLeft()

		colorIndex := i % len(colors)

		for _, line := range lines {
			out := termenv.String(line)
			fmt.Println(out.Foreground(colors[colorIndex]).String())
		}
		time.Sleep(150 * time.Millisecond)
	}
	s.Stop()
}

func main() {

	// Set window title and size
	setWindowTitleAndSize("KREX NEST.RIP BOOSTER v4.0 - ULTIMATE EDITION")

	rand.Seed(time.Now().UnixNano())

	screen.Clear()
	screen.MoveTopLeft()
	animateLogo()

	screen.Clear()
	screen.MoveTopLeft()
	titleColor.Print(logo)
	fmt.Println()

	infoColor.Println("═══════════════════════════════════════════════════════════════")
	subtitleColor.Println("      NEST.RIP TRAFFIC BOOSTER - ULTIMATE EDITION v4.0     ")
	infoColor.Println("═══════════════════════════════════════════════════════════════")

	warningColor.Println("🚀 The most powerful nest.rip booster available")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	infoColor.Print("Enter the nest.rip link (format: ")
	urlColor.Print("https://nest.rip/s/hQzxp")
	infoColor.Println("):")
	highlightColor.Print("➤ ")

	linkInput, _ := reader.ReadString('\n')
	url := strings.TrimSpace(linkInput)

	if !strings.HasPrefix(url, "https://nest.rip/") {
		errorColor.Println("❌ Invalid URL format. URL must start with https://nest.rip/")
		return
	}

	fmt.Println()
	infoColor.Println("Select click type:")
	highlightColor.Println("1. ⚡ Unique Clicks (appears as different visitors)")
	highlightColor.Println("2. 🔥 Total Clicks (faster, but may count as same visitor)")
	highlightColor.Print("➤ ")

	clickTypeInput, _ := reader.ReadString('\n')
	clickType := strings.TrimSpace(clickTypeInput)

	isUniqueClick := false
	if clickType == "1" {
		isUniqueClick = true
		successColor.Println("\n✅ Selected: Unique Clicks")
	} else if clickType == "2" {
		successColor.Println("\n✅ Selected: Total Clicks")
	} else {
		warningColor.Println("\n⚠️ Invalid option. Defaulting to Total Clicks.")
	}

	infoColor.Println("\nEnter the number of unique clicks to generate:")
	highlightColor.Print("➤ ")

	clicksInput, _ := reader.ReadString('\n')
	clicks, err := strconv.Atoi(strings.TrimSpace(clicksInput))
	if err != nil || clicks <= 0 {
		errorColor.Println("❌ Invalid number of clicks. Please enter a positive integer.")
		return
	}

	totalClicks := clicks
	if isUniqueClick {
		infoColor.Println("\nEnter the total number of clicks to generate (should be higher than unique clicks):")
		highlightColor.Print("➤ ")

		totalClicksInput, _ := reader.ReadString('\n')
		totalClicksVal, err := strconv.Atoi(strings.TrimSpace(totalClicksInput))
		if err != nil || totalClicksVal < clicks {
			warningColor.Printf("⚠️ Invalid input. Setting total clicks to %d\n", clicks)
		} else {
			totalClicks = totalClicksVal
		}
	}

	screen.Clear()
	screen.MoveTopLeft()
	titleColor.Print(logo)
	fmt.Println()

	infoColor.Println("═══════════════════════════════════════════════════════════════")
	infoColor.Print("🚀 Starting to generate ")
	numberColor.Printf("EXACTLY %d", totalClicks)
	infoColor.Print(" clicks")

	if isUniqueClick {
		infoColor.Printf(" (%d unique, %d total)", clicks, totalClicks)
	}

	infoColor.Print(" for ")
	urlColor.Println(url)
	infoColor.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	var successCount int32
	startTime := time.Now()

	remainingClicks := int32(totalClicks)

	maxConcurrent := 150
	workerQueue := make(chan int, maxConcurrent)
	results := make(chan RequestResult, maxConcurrent)
	done := make(chan bool)

	s := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
	s.Prefix = " "
	s.Start()

	ticker := time.NewTicker(200 * time.Millisecond)
	go func() {
		lastUpdate := time.Now()
		lastCount := int32(0)

		for {
			select {
			case <-ticker.C:
				current := atomic.LoadInt32(&successCount)
				now := time.Now()
				elapsed := now.Sub(lastUpdate)
				rate := float64(current-lastCount) / elapsed.Seconds()
				if rate < 0 {
					rate = 0
				}

				fmt.Print("\r\033[K")

				progressBar := displayProgressBar(current, int32(totalClicks), 40)
				progressPercent := float64(current) / float64(totalClicks) * 100

				infoColor.Printf("\r%s ", progressBar)
				numberColor.Printf("%.1f%%", progressPercent)
				infoColor.Printf(" | ")
				successColor.Printf("%d/%d", current, totalClicks)
				infoColor.Printf(" clicks | ")
				highlightColor.Printf("%.2f", rate)
				infoColor.Printf(" clicks/sec | ETA: ")
				warningColor.Printf("%s", calculateETA(current, int32(totalClicks), rate))

				lastCount = current
				lastUpdate = now
			case <-done:
				return
			}
		}
	}()

	totalToLaunch := totalClicks + 100
	go func() {
		for i := 0; i < totalToLaunch; i++ {
			workerQueue <- i
		}
	}()

	var wg sync.WaitGroup

	for i := 0; i < maxConcurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for workerID := range workerQueue {
				if atomic.LoadInt32(&remainingClicks) <= 0 {

					break
				}

				if workerID%10 == 0 {
					time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
				}

				clickWorker(url, results, isUniqueClick, workerID, clicks, &remainingClicks)
			}
		}()
	}

	go func() {
		for result := range results {
			if result.Success {
				atomic.AddInt32(&successCount, 1)
			}

			if atomic.LoadInt32(&successCount) >= int32(totalClicks) {

				atomic.StoreInt32(&remainingClicks, 0)

				close(workerQueue)
				break
			}
		}
	}()

	wg.Wait()
	close(results)
	done <- true
	ticker.Stop()
	s.Stop()

	elapsedTime := time.Since(startTime)

	screen.Clear()
	screen.MoveTopLeft()

	titleColor.Print(logo)
	fmt.Println()

	infoColor.Println("═══════════════════════════════════════════════════════════════")
	successColor.Printf("✅ MISSION ACCOMPLISHED! Generated EXACTLY %d clicks\n", totalClicks)

	if isUniqueClick {
		highlightColor.Printf("🔹 Unique Clicks: %d\n", clicks)
		highlightColor.Printf("🔹 Total Clicks: %d\n", totalClicks)
	}

	highlightColor.Printf("⏱️ Time taken: %s\n", elapsedTime)
	highlightColor.Printf("⚡ Average speed: %.2f clicks per second\n", float64(totalClicks)/elapsedTime.Seconds())

	speedRating := "🚀 LUDICROUS SPEED"
	if float64(totalClicks)/elapsedTime.Seconds() > 300 {
		speedRating = "🌌 INTERSTELLAR SPEED"
	} else if float64(totalClicks)/elapsedTime.Seconds() < 50 {
		speedRating = "🐢 DECENT SPEED"
	}

	successColor.Printf("🏆 Performance rating: %s\n", speedRating)
	infoColor.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()
	subtitleColor.Println(copyright)
}

func calculateETA(current, total int32, rate float64) string {
	if rate <= 0 || current >= total {
		return "Done"
	}

	remaining := float64(total - current)
	seconds := remaining / rate

	if seconds < 60 {
		return fmt.Sprintf("%.0fs", seconds)
	} else if seconds < 3600 {
		return fmt.Sprintf("%.1fm", seconds/60)
	} else {
		return fmt.Sprintf("%.1fh", seconds/3600)
	}
}
