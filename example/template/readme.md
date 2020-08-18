# snippet

Snippet finds code snippets delimited by comments in source code and
expands templates that use these snippets.

The core use-case is to enable developers to test code snippets that they put into documentation in CI.

## Installation

```shell script
go install github.com/snamber/snippet/cmd/snippet
```

## Usage

`snippet` offers four command line parameters

```text
-h --help        Displays help with available flag, subcommand, and positional value parameters.
-s --snippets    The file or directory in which snippets are defined.
-t --templates   The file or directory containing the output template(s).
-o --output      The directory where output file(s) will be written.
```

To generate this readme for example, clone this project, and execute the following command in the root directory

```shell script
snippet -s example/testdata -t example/template/readme.md -o . 
```
or 
```shell script
snippet -s example/testdata -t example/template -o . 
```

### How to define a snippet

Single-line comments like
```go
// START SNIPPET tag
```
define the start of snippets, and associate the snippet with a tag.

Single-line comments like
```go
// END SNIPPET
```
define the end of the respective snippet. 

The prefix used to declare a comment depends on the language; in Python a snippet would be defined like this:
```python
# START SNIPPET tag
# END SNIPPET
```

Tags can contain letters, digits, `_` and `-`. Spacing matters.

### How to use a snippet

To mark a location into which a snippet should be expanded, write a line like
```md
// PUT SNIPPET tag　
```
Reference the respective tag. Markdown does not have comments: Use `//` as a comment prefix in Markdown.

### Indentation

Snippet trims the indentation at the _definition_ of a snippet, and replaces it with the indentation of the 'PUT' comment at the _use_
of a snippet.

For the following definition
```go
    // START SNIPPET tag
    func foo() {
        var foo string
    }   
    // END SNIPPET
```
and the following use
```md
// PUT SNIPPET tag　
```
would result in this expanded snippet
```go
func foo() {
    var foo string
}   
```

## Supported languages

Currently snippet supports the following file extensions, with the associated single-line comment prefixes

```text
.c:    "//"
.cpp:  "//"
.go:   "//"
.h:    "//"
.java: "//"
.js:   "//"
.kt:   "//"
.md:   "//"
.py:   "#"
.ts:   "//"
.sh:   "//"
.txt:  "//"
```

## How to contribute

Contributions, specifically support for more languages, is appreciated.

To support another file extension modify language.go, provide a snippet for the readme by putting a file
of the respective language into `example/testdata`, updating `example/template`, re-generate the readme by using `snippet`, and create a PR. 

## Examples

The folder `test` contains an example where a snippet extracts part of a tested piece of code.

Here we see sample output of `snippet`, with snippets defined in `example/testdata`, and the template for this file defined in `example/template`.

C
```c
// PUT SNIPPET c
```

C++
```cpp
// PUT SNIPPET cpp
```

Go
```go
// PUT SNIPPET go
```

Java
```java
// PUT SNIPPET java
```

Javascript
```js
// PUT SNIPPET javascript
```

Kotlin
```kotlin
// PUT SNIPPET kotlin
```

Python
```python
// PUT SNIPPET python
```

Text
```text
// PUT SNIPPET txt
```

Typescript
```typescript
// PUT SNIPPET typescript
```
