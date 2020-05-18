# Streaming server
A multicast server using Nginx.

## Prerequisites

- [Vagrant](https://www.vagrantup.com/)
- [Virtualbox (or other Vagrant supported provider)](https://www.virtualbox.org/)
- [Ansible (if not on Windows)](https://www.ansible.com/)

### Installation
To install the prerequisites above you can navigate to each link above and follow the instructions. Alternatively, you can use a package manager like [Homebrew](https://brew.sh/) or [Chocolatey](https://chocolatey.org/). Here are some quick scripts to install the prerequisites above with either of these two package managers.

Homebrew:
```sh
brew cask install \
    virtualbox \
    vagrant &&
brew install ansible
```

Chocolatey:
```powershell
cinst virtualbox vagrant -y
```

## Configuration

To configure the server, you will need to fetch RTMP urls and stream keys for the services you want to stream to and add them to the list of keys in the `bootstrap.yml` file. For example, if you'd like to stream to Facebook and YouTube at the same time, find the URLs and stream keys for both and replace this:

```yaml
streams:
    - url: "webcast.sermonaudio.com/sa"
      key: "some-sa-key"
   
```

With This:

```yaml
streams:
    - url: "rtmp.facebook.com/live1"
      key: "some-facebook-key"
    - url: "a.rtmp.youtube.com/live2"
      key: "some-youtube-key"
```

You can also add any number of other RTMP stream provider endpoints by just adding another key in the map.

```yaml
streams:
    - url: "rtmp.facebook.com/live1"
      key: "some-facebook-key"
    - url: "a.rtmp.youtube.com/live2"
      key: "some-youtube-key"
    - url: "rtmp.twitch.com/live3"
      key: "some-twitch-key"
    - url: "rtmp.example.com/live4"
      key: "some-example-key"
```

All of the keys above are example keys and don't reflect what the URLs or keys will actually be. Look at your stream providers documentation to get the URL and key you'll need.

## Installation

Clone the repository, make your config changes, and run `vagrant up` in your terminal. The scripts should run and build out your server automatically!

## Using the server

Open up your OBS streaming settings and point your URL to `rtmp://10.0.1.15/live`, type in something for your stream key. It doesn't matter what your local stream key is. Start your stream and verify that the services you configured are receiving the stream.

### Sources:
- [OBS Forum](http://bit.ly/2NTRGSm)
- [Video Guide](https://www.youtube.com/watch?v=o6N_fX5IcLM)