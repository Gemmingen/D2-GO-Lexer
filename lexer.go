package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

//Hier den Dateiname eingeben :

const filename = "D2/exampleuml.d2"

// Console -> go run lexer.go

//Das Script erstellt eine "Output.go" mit den Interfaces und Structs aus der D2. Ebenso werden Methoden und Variablen übernommen. Methoden übernehmen Parameter- und Returnvalues.
//Verbindungspfeile in D2 z.B. "Interface1 -> Interface2" müssen ausgeblendet werden, durch "#STOP".
//Das Skript hört dann auf die D2 zu lesen, wenn es "#STOP" erkennt und fängt an zu schreiben.

type Class struct {
	Name       string
	Properties []Property
	Methods    []Method
}

type Property struct {
	Name string
	Type string
}

type Method struct {
	Name   string
	Params string
	Return string
}

// Hilfsfunktionen

func keepOnlyLetters(input string) string {
	re := regexp.MustCompile(`[^a-zA-Z]+`)
	return re.ReplaceAllString(input, "")
}

// Funktion zur Generierung von Go-Interfaces basierend auf den Klassen
func generateGoInterfaces(classes []*Class, writer *os.File) {
	fmt.Fprintln(writer, "// Generated Go Interfaces")
	fmt.Fprintln(writer, "package main")

	for _, class := range classes {
		// Wenn keine Methoden vorhanden sind, dann ein Struct generieren
		if len(class.Methods) == 0 {
			fmt.Fprintf(writer, "type %s struct {\n", class.Name)
			for _, property := range class.Properties {
				fmt.Fprintf(writer, "  %s %s\n", property.Name, property.Type)
			}
			fmt.Fprintln(writer, "}")
		}

		// Wenn Methoden vorhanden sind, dann ein Interface generieren
		if len(class.Methods) > 0 {
			fmt.Fprintf(writer, "\ntype %s interface {\n", class.Name)
			for _, method := range class.Methods {
				fmt.Fprintf(writer, "  %s(%s) %s\n", method.Name, method.Params, method.Return)
			}
			fmt.Fprintln(writer, "}")
			fmt.Println("Go interfaces generated successfully!")
		}
	}
}

func createClass(Lines []string, startLine int) ([]*Class, int) {
	Linecounter := startLine
	var classes []*Class
	var currentClass *Class

	for Linecounter < len(Lines) && !strings.Contains(Lines[Linecounter], "#STOP") {
		line := strings.TrimSpace(Lines[Linecounter])

		// Neue Klasse erkannt
		if strings.Contains(line, "{") {
			// Füge vorherige Klasse hinzu (falls vorhanden)
			if currentClass != nil {
				classes = append(classes, currentClass)
			}
			className := keepOnlyLetters(line)
			currentClass = &Class{Name: className}
			fmt.Printf("Found class: %s\n", className)
		} else if currentClass != nil {
			// Eigenschaften erkennen (Format: name: type)
			if !strings.Contains(line, "shape") {
				if strings.Contains(line, ":") && !strings.Contains(line, "(") {
					lineSplit := strings.Split(line, ":")
					if len(lineSplit) == 2 {
						property := Property{
							Name: strings.TrimSpace(lineSplit[0]),
							Type: strings.TrimSpace(lineSplit[1]),
						}
						currentClass.Properties = append(currentClass.Properties, property)
						fmt.Printf("Added property: %s : %s\n", property.Name, property.Type)
					}
				}

				// Methoden erkennen (Format: methodName(params): returnType)
				if strings.Contains(line, "(") {
					lineSplit := strings.Split(line, ":")
					if len(lineSplit) == 2 {
						methodName := strings.TrimSpace(lineSplit[0])
						returnType := strings.TrimSpace(lineSplit[1])

						// Parameter extrahieren
						re := regexp.MustCompile(`\((.*?)\)`)
						match := re.FindStringSubmatch(methodName)
						params := ""
						if match != nil {
							params = match[1]
						}

						method := Method{
							Name:   strings.TrimSpace(strings.Split(methodName, "(")[0]),
							Params: params,
							Return: returnType,
						}
						currentClass.Methods = append(currentClass.Methods, method)
						fmt.Printf("Added method: %s(%s) : %s\n", method.Name, method.Params, method.Return)
					}
				}
			}
		}

		Linecounter++
	}

	// Letzte Klasse hinzufügen
	if currentClass != nil {
		classes = append(classes, currentClass)
	}
	return classes, Linecounter
}

func FileToSlice(file string) []string {
	content, err := os.Open(file)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}
	var Lines []string
	scanner := bufio.NewScanner(content)

	for scanner.Scan() {
		Lines = append(Lines, scanner.Text())
	}
	return Lines
}
func main() {
	outputFile := "output.go"
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	defer file.Close()

	// Datei einlesen

	Lines := FileToSlice(filename)

	// Alle Klassen sammeln
	var allClasses []*Class
	Linecounter := 0

	for Linecounter < len(Lines) {
		var classes []*Class
		classes, Linecounter = createClass(Lines, Linecounter)
		allClasses = append(allClasses, classes...)
		Linecounter++
	}

	// Go-Interfaces generieren
	generateGoInterfaces(allClasses, file)

}
