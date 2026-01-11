# Youtube transcriber
This is an orchestrator app.

Combine this with the other apps.
- [ffmpeg](https://www.ffmpeg.org/download.html)
- [yt-dlp](https://github.com/yt-dlp/yt-dlp/wiki/Installation)
- [whisper.cpp](https://github.com/ggml-org/whisper.cpp)

before use app prepare and install apps below.


## Quick start
### Whisper.cpp
First make separate directory for clone the repository:

```bash
git clone https://github.com/ggml-org/whisper.cpp.git
```

Navigate into the directory:

```
cd whisper.cpp
```

Then, download one of the Whisper [models](models/README.md) converted in [`ggml` format](#ggml-format). For example:

```bash
sh ./models/download-ggml-model.sh base.en
```
If need, you can update `download-ggml-model.sh` from [download-ggml-model.sh](https://raw.githubusercontent.com/ggml-org/whisper.cpp/refs/heads/master/models/download-ggml-model.sh)

Now build the [whisper-cli](examples/cli) example and transcribe an audio file like this:

```bash
# build the project
cmake -B build
cmake --build build -j --config Release
```
You can com result `./build/bin/whisper-cli` to the project path `<project dir>/ext-tools/whisper/whisper-cli`
```bash
# transcribe an audio file
cp ./build/bin/whisper-cli <project dir>/ext-tools/whisper/whisper-cli
```
## ytdlp and ffmpeg
You must install by instructions.
- [yt-dlp](https://github.com/yt-dlp/yt-dlp/wiki/Installation)
- [whisper.cpp](https://github.com/ggml-org/whisper.cpp)

## Build

```shell
make build

chmod +x ./ytb
```

## Run

```shell
./ytb
```

follow the instructions in UI.

Enjoy!