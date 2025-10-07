package tui

import (

	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"codeline/llm"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/term"
)

const (
	userBubble   = "[#34e2e2]"
	ollamaBubble = "[#ad7fa8]"
	resetColor   = "[white]"
)

func formatBubble(text, color string, width int) string {
	return fmt.Sprintf("%s%s%s", color, text, resetColor)
}
func StartChat(ctx context.Context, client llm.LLM) {
	app := tview.NewApplication()

	chatView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() { app.Draw() })
	chatView.SetBorder(true)
	chatView.SetBackgroundColor(tcell.GetColor("#1a1a1a"))

	label := tview.NewTextView().
		SetText(">").
		SetTextColor(tcell.ColorWhite).
		SetBackgroundColor(tcell.GetColor("#0f0f0f"))

	inputField := tview.NewInputField().
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.GetColor("#0f0f0f")).
		SetFieldTextColor(tcell.ColorWhite)
	inputField.SetBackgroundColor(tcell.GetColor("#0f0f0f"))

	messages := []string{}

	introText := `
===========================================================
===     =============  =========  =========================
==  ===  ============  =========  =========================
=  ==================  =========  =========================
=  =========   ======  ===   ===  ========  ==  = ====   ==
=  ========     ===    ==  =  ==  ============     ==  =  =
=  ========  =  ==  =  ==     ==  ========  ==  =  ==     =
=  ========  =  ==  =  ==  =====  ========  ==  =  ==  ====
==  ===  ==  =  ==  =  ==  =  ==  ========  ==  =  ==  =  =
===     ====   ====    ===   ===        ==  ==  =  ===   ==
===========================================================

Welcome to the chat interface!
Type your message and press Enter.
Type "exit", "quit" or "q" to quit.
`
	messages = append(messages, introText)
	updateChat(chatView, messages)

	firstMessageSent := false

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			return
		}

		input := strings.TrimSpace(inputField.GetText())
		if input == "" {
			return
		}
		if strings.EqualFold(input, "exit") || strings.EqualFold(input, "quit") || strings.EqualFold(input, "q") {
			app.Stop()
			return
		}

		width, _, err := term.GetSize(int(os.Stdin.Fd()))
		if err != nil {
			width = 80
		}
		bubbleWidth := int(float64(width) * 0.6)

		// Remove intro after first message
		if !firstMessageSent {
			firstMessageSent = true
			messages = messages[1:]
		}

		// User bubble
		messages = append(messages, formatBubble(input, userBubble, bubbleWidth))
		updateChat(chatView, messages)

		// Placeholder for the assistant’s streamed message
		streamBuffer := new(strings.Builder)
		messages = append(messages, formatBubble("", ollamaBubble, bubbleWidth))
		updateChat(chatView, messages)

		// Start streaming response in a goroutine
		go func() {
			stream, err := client.AskStream(ctx, input)
			if err != nil {
				app.QueueUpdateDraw(func() {
					messages[len(messages)-1] = formatBubble("Error: "+err.Error(), ollamaBubble, bubbleWidth)
					updateChat(chatView, messages)
				})
				return
			}

			for token := range stream {
				streamBuffer.WriteString(token)

				// Update UI safely
				app.QueueUpdateDraw(func() {
					messages[len(messages)-1] = formatBubble(streamBuffer.String(), ollamaBubble, bubbleWidth)
					updateChat(chatView, messages)
				})
			}
		}()

		inputField.SetText("")
	})

	inputFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(label, 1, 0, false).
		AddItem(inputField, 0, 1, true)
	inputFlex.SetBackgroundColor(tcell.GetColor("#0f0f0f"))

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexColumn).
				AddItem(nil, 2, 0, false).
				AddItem(chatView, 0, 1, false).
				AddItem(nil, 1, 0, false),
			0, 1, false,
		).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexColumn).
				AddItem(nil, 2, 0, false).
				AddItem(inputFlex, 0, 1, true).
				AddItem(nil, 1, 0, false),
			2, 0, true,
		).
		AddItem(nil, 1, 0, false)

	app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		screen.Fill(' ', tcell.StyleDefault.Background(tcell.GetColor("#1a1a1a")))
		return false
	})

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		log.Fatalf("Error starting chat: %v", err)
	}
}

func updateChat(view *tview.TextView, messages []string) {
	view.Clear()
	for _, msg := range messages {
		fmt.Fprintln(view, msg)
	}
}
