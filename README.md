Einlesen von UML-Daten in D2 Format: 

Das Script erstellt eine "Output.go" mit den Interfaces und Structs aus der D2. Ebenso werden Methoden und Variablen übernommen. Methoden übernehmen Parameter- und Returnvalues.
Verbindungspfeile in D2 z.B. "Interface1 -> Interface2" müssen ausgeblendet werden, durch "#STOP".
Das Skript hört dann auf die D2 zu lesen, wenn es "#STOP" erkennt und fängt an zu schreiben.

-Das Skript hört dann auf die D2 zu lesen, wenn es "#STOP" erkennt und fängt an zu schreiben.

-Erkennung von Klassen: Die Klassen werden durch das Erkennen von '{}' Blöcken identifiziert.

-Erkennung von Eigenschaften: Eigenschaften werden im Format 'name: type' erkannt.

-Erkennung von Methoden: Methoden werden im Format 'methodName(params): returnType' erkannt.

-Generierung von Go-Interfaces: Für jede erkannte Klasse wird ein Go-Interface erstellt. 

Bevor du jetzt deine UML Umwandel kannst, musst du sie in den D2 Ordner schieben und in lexer.go die const filename anpassen. 

Console -> go run lexer.go

https://d2lang.com/tour/uml-classes 


//Hallo yt :D
