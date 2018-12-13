// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START translate_quickstart]
// Sample translate-quickstart translates "Hello, world!" into Russian.
package trans

import (
	"fmt"
	"log"
	"spos_lang/ulog"

	// Imports the Google Cloud Translate client package.
	"cloud.google.com/go/translate"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
	"html"
	//	"strings"
)

func TransRequest(s_ln string, t_ln string, text string) (string, error) {

	if s_ln != t_ln {
		ctx := context.Background()

		// Creates a client.
		client, err := translate.NewClient(ctx)
		if err != nil {
			fmt.Printf("%s: Trans failed to create client: %v\n", ulog.CreateDateString(), err)
			return "", err
		}

		// Sets the text to translate.
		//	text := "Hello, world!"
		// Sets the target language.
		target, err := language.Parse(t_ln)
		if err != nil {
			fmt.Printf("%s: Trans failed to parse target language:  %v\n", ulog.CreateDateString(), err)
			return "", err
		}

		// Translates the text into Russian.
		opt := translate.Options{
			Source: language.Make(s_ln),
		}

		translations, err := client.Translate(ctx, []string{text}, target, &opt)
		if err != nil {
			fmt.Printf("%s: Trans failed to translate text:  %v\n", ulog.CreateDateString(), err)
			return "", err
		}

		return html.UnescapeString(translations[0].Text), nil
	} else {
		return "", nil
	}

}

func trans_main() {
	ctx := context.Background()

	// Creates a client.
	client, err := translate.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the text to translate.
	text := "Hello, world!"
	// Sets the target language.
	target, err := language.Parse("ru")
	if err != nil {
		log.Fatalf("Failed to parse target language: %v", err)
	}

	// Translates the text into Russian.
	translations, err := client.Translate(ctx, []string{text}, target, nil)
	if err != nil {
		log.Fatalf("Failed to translate text: %v", err)
	}

	fmt.Printf("Text: %v\n", text)
	fmt.Printf("Translation: %v\n", translations[0].Text)
}

// [END translate_quickstart]
