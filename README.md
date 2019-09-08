[![CircleCI](https://circleci.com/gh/vatolvan/go-hue/tree/master.svg?style=svg)](https://circleci.com/gh/vatolvan/go-hue/tree/master)

# go-hue
Hue light switcher with Go

## Config

Create `config.json` to the root of the project with the following

```{
  "hue_bridge_username": "<username you have created to hue bridge>",
  "hue_bridge_ip": "<ip of your hue bridge>"
}
```


Check out [Hue developer guide](https://developers.meethue.com/develop/get-started-2/) for instructions on how to get `username` and the `ip`.