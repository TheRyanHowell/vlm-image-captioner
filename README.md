# VLM Image Captioner

`vlm-image-captioner` is a command-line tool that generates captions for one or more images using a Vision Language
Model. It can output the captions to standard output or in CSV format.

## Installation

You can install it using `go install`:

```bash
go install github.com/theryanhowell/vlm-image-captioner@latest
```

Alternatively, you can build the program with:

```bash
go build .
```

## Usage

To generate a caption for a single image:

```bash
vlm-image-captioner /path/to/your/image.jpg
```

To generate captions for multiple images:

```bash
vlm-image-captioner /path/to/image1.jpg /path/to/image2.png
```

By default, the output will be printed to standard output.

### CSV Output

To output the captions in CSV format, use the `--csv` or `-c` flag:

```bash
vlm-image-captioner --csv /path/to/image1.jpg /path/to/image2.png > captions.csv
```

This will create a `captions.csv` file with the following columns: `imagepath` and `caption`.

## Configuration

The tool is configured using the following environment variables:

* `OPENAI_API_KEY`: Your OpenAI API key. This is a required variable when using the official OpenAI API.
* `OPENAI_BASE_URL`: The base URL for the OpenAI API. This is optional and defaults to the official OpenAI API URL.
* `OPENAI_MODEL`: The model to use for captioning. This is optional and defaults to `gpt-5`.

### Local VLM Usage

You can use the above environment variables to run against any OpenAI compatible API, such as
with [LM Studio](https://lmstudio.ai).

For example, the following command uses the local LM Studio API and Google's Gemma 3 4B Vision Language Model.

```bash
OPENAI_BASE_URL="google/gemma-3-4b" OPENAI_BASE_URL="http://localhost:1234/v1" vlm-image-captioner images/*
```
