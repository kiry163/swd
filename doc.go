// Package swd provides a sensitive word detection engine.
//
// Call New to create an engine, load user-defined words with loaders or direct
// update methods, then use Contains, Find, FindAll, or Replace to inspect text.
// Words, Export, and ExportFile expose the current word table as a snapshot or
// as the same simple text format accepted by the built-in loaders. Matching can
// normalize case, width, Traditional Chinese variants, symbols, and common
// low-risk similar-character substitutions according to Options.
package swd
