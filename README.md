# Fusion #

Is a javascript bundling plugin designed to work like sass. Configure it at startup. Re-bundle when you see a new request.

## How? ##

### Configuration ###

-  You can specify multiple bundles in one yaml file.
-  Per bundle, specify 
  +  The output file name
  +  An input file list (whose order is preserved)
  +  An input directory

The unique set of files specified this way will be combined into one file. 

### Modes ###

The quick mode just does a dumb concatentation. This is very fast. Its recommended that you use this mode when using the 'reloading' feature.

The optimized mode uses Google Closure Compiler's SIMPLE_OPTIMIZATIONS flag -- which means comments and whitespace are stripped, and basic (non-obfuscating) optimizations are performed (such as inlining a function thats only called once).

### Example ###

See the example [bundle.yaml](doc/example-bundles.yaml) file. This example config will create 3 bundles: main.js / checkout.js / bottom.js