package config

type fileConfig struct {
	TextOutput string `mapstructure:"textoutput"`

	Youtube struct {
		YtDlPath string `mapstructure:"ytdlpath"`
	} `mapstructure:"youtube"`

	Ffmpeg struct {
		Path string `mapstructure:"path"`
	} `mapstructure:"ffmpeg"`

	Transcriber struct {
		Language string `mapstructure:"language"`

		Whisper struct {
			Cli   string `mapstructure:"cli"`
			Model string `mapstructure:"model"`
		} `mapstructure:"whisper"`

		OpenAI struct {
			Model string `mapstructure:"model"`
		} `mapstructure:"openai"`
	} `mapstructure:"transcriber"`
}
