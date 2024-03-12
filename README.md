# Lanty server component

[![Go](https://github.com/seternate/go-lanty/actions/workflows/go.yaml/badge.svg)](https://github.com/seternate/go-lanty/actions/workflows/go.yaml)

Hey hey, gamers! 🎮👾

So, guess what? There's this awesome new app called Lanty, and let me tell ya, it's like the ultimate tool for LAN parties! Picture this: you're at your buddy's place, everyone's hyped to game together, and boom! Lanty swoops in to save the day!

First off, Lanty hooks you up with all the games you need. No more hunting around for CDs or downloads. It's all right there, ready to roll.

And get this – setting up game servers? Piece of cake! With Lanty, it's like hitting the "easy" button. Just a few clicks, and boom, you're hosting your own server like a pro.

Now, remember how sometimes it's a pain trying to figure out who's who at a LAN party? Not anymore! Lanty's got your back with easy-peasy participant discovery. Plus, it automatically picks up and shares IP addresses, so no more hassle!

So, what are you waiting for? Grab Lanty, gather your crew, and let the LAN party madness begin! 🚀🎉

## 🚀 Ready to Get Lanty Up and Running? Let's Do This! 🚀

Alright, fellow gamers, it's time to dive into the installation process for Lanty! Here's what you gotta do:

#### Step 1

Head on over to the [latest release](https://www.github.com/seternate/go-lanty/releases) of the server application and grab that bad boy. Oh, and don't forget to snag the [client application](https://www.github.com/seternate/go-lanty-client/releases) too, both with the same version, of course. We need 'em to be on the same page! 🦖

#### Step 2

Once you've got those files, it's time to work some magic. Extract the server artifact like a pro and then slide that client zip artifact into the `download` directory of the server application. Easy peasy, right?

🎉 __Woohoo! You did it!__ 🎉

Now give yourself a pat on the back because you, my friend, are now officially a Lanty pro! Get ready to level up your LAN party game and dive into some serious gaming fun. 🎮👾

## 🌱 Let's Get Lanty All Set Up and Ready to Roll! 🌱

Before Lanty can start showcasing all those epic games to your clients, we've got a bit of setup to do. But don't worry, it's nothing you can't handle! Here's the lowdown:

### 📦 Add game data 📦

Alrighty, first things first. We gotta get those game files zipped up and ready to roll. Pop each game files into a zip archive and toss the bad boys into the `game-data` folder. Oh, and remember to name the zip file with a `slug` so the server can find it later. (For example, if it's __Call of Duty 4__, you'd name it __call-of-duty-4.zip__)

❗ Pro Tip: Keep it clean! Zip those game files without any extra subfolders. We want 'em neat and tidy. ❗

    ✅                          ⛔
    ├── call-of-duty-4.zip      ├── call-of-duty-4.zip
    │   ├── Docs/               │   ├── subfolder/
    │   │   ├── ...             │   │   ├── Docs/
    |   ├── Mods/               |   |   |   ├── ...
    |   |   ├── ...             |   |   ├── Mods/
    │   ├── iw3mp.exe           │   |   |   ├── ...
    │   ├── ...                 │   |   ├── iw3mp.exe
                                |   |   ├── ...

### 🖼️ Add game icon 🖼️

Next up, let's give those games some flair! Every game needs its own icon, right? Make sure it's a .png file, or else it won't get picked up. Drop those icons into the `game-icon` folder and rename 'em using the same `slug` as the game data zip files. (For example, if it's __Call of Duty 4__, you'd name it __call-of-duty-4.png__)

### ⚙️ Add game config ⚙️

Last but certainly not least, we need to add a configuration file for each game. These files are written using yaml, so get ready to flex those config muscles! Need help? Check out the [game configuration](#game-configuration-files) files section for all the deets.

And there you have it, folks! Lanty is well on its way to becoming your ultimate LAN party wingman. Now go forth, add those games, and let the gaming adventures begin! 🎮🚀

🎉 __Woohoo! You are going above and beyond!__ 🎉

And there you have it, fellow gamer! Lanty is poised and ready to elevate your LAN party experience to legendary status. So what are you waiting for? Dive in, unleash those gaming masterpieces, and let the fun times roll! Remember, with Lanty by your side, the gaming world is your oyster. 🌟🎮 Let's game on! 🚀🔥

## 🎮 Ready to Get the Party Started with Lanty? Gogogo! 🎮

Alrighty, once you've wrapped up the setup for Lanty, it's time to dive right into the action! Here's what you gotta do:

__🏃Run Lanty 🏃__

It's as simple as that – just fire up Lanty and you're good to go! No fuss, no muss, just pure gaming goodness awaits you.

### ☁️ Download the Client ☁️

Now, here's where the magic happens! Users can easily grab the client application by cruising over to `http://<server-url>:<port>/download` with their trusty web browser. A quick download will kick off, and all they gotta do is unzip and they're basically ready to roll! For more juicy details on how to use the client application like a pro, check out its [README](https://www.github.com/seternate/go-lanty).

#### 🎮 Let the LANParty Begin! 🎮

With Lanty in your corner, it's time to crank up the excitement and let the LANParty madness commence! Gather your pals, fire up those games, and get ready for an epic gaming extravaganza like no other! Let's make some gaming memories that'll last a lifetime! 🚀🔥

## 🛠️ Game Configuration Files 🛠️

Alright, let's get down to business with the nitty-gritty details of game configuration files for Lanty. This might sound a tad boring, but hey, it's essential stuff! Here's what you need to know.

### Required Fields - Basic Configuration

At its core, a game configuration file is a yaml file containing a few essential fields:

* Required
    * [`slug`](#slug) - A unique identifier for the game
    * [`name`](#name)- A human-readable name for the game
    * [`client`](#client) - Game startup information for launching the game and joining a server
        * [`executable`](#client) - Relative path to the game's executable within the game folder

* Optional
    * [`client`](#client)
        * [`arguments`](#specifying-arguments) - Command-line options for the game executable
            * [`seperator`](#specifying-seperator) - Specifies how command-line options are separated and parsed
            * [`items`](#specifying-items) -  A list of individual command-line options
    * [`server`](#server) - Game startup information for starting a server
        * [`executable`](#server) - Relative path to the game's server executable within the game folder
        * [`arguments`](#specifying-arguments) - Command-line options for the server executable
            * [`seperator`](#specifying-seperator) - Specifies how command-line options are separated and parsed
            * [`items`](#specifying-items) -  A list of individual command-line options

### 🔍 Simple things first 🔍

Alright, let's simplify things a bit. If a game doesn't have any command-line options available, then all you need is the following configuration file:

```yaml
slug: "call-of-duty-4"
name: "Call of Duty 4"
client:
  executable: "iw3mp.exe"
```

#### `slug`

The slug serves as a unique identifier for each game. It should only contain lowercase characters, numerals, and dashes. This identifier is crucial for locating the game data zip file and game icon file.

#### `name`

The name field is a user-friendly string used to display the game's name to users.

#### `client`

The client section specifies the relative path within the game folder to the game's `executable`, used for launching the game and joining a server. Additionally, it can include command-line arguments for the executable. For more details on specifying arguments, check out the [arguments section](#specifying-arguments).

```yaml
slug: "call-of-duty-4"
name: "Call of Duty 4"
client:
  executable: "iw3mp.exe"
  arguments:
    items:
    - ...
```

#### `server`

Now, let's talk about the server section. This part of the configuration file specifies the relative path within the game folder to the `executable` used for starting a server. It also includes any command-line arguments for the executable. For more detailed information, check out the [arguments section](#specifying-arguments).

```yaml
slug: "call-of-duty-4"
name: "Call of Duty 4 - Modern Warfare"
client:
  executable: "iw3mp.exe"
server:
  executable: "iw3mp.exe"
  arguments:
    items:
    - ...
```

Alright, there you have it – the essential rundown on game configuration files for Lanty. It might not be the most thrilling read, but mastering these details will ensure smooth sailing in your gaming adventures.

### 🕵️‍♂️🚀 Embarking on the Adventure of Configuration 🚀🕵️‍♂️

Hold on tight, fellow explorer! We're diving deep into the realm of configuration, where only the bravest dare to tread. This section is for the experts, the trailblazers ready to conquer command-line complexities and emerge victorious.

Are you ready? Let's do this!

### Specifying `arguments`

Welcome, fellow adventurer, to the realm of arguments – a crucial aspect of command-line configuration. Here, we wield the power to specify command-line arguments for an executable, shaping its behavior with precision and finesse.

But what exactly does this entail? Let's break it down:

* `arguments` serves as our gateway to customizing the behavior of executables.
* It comprises a __global__ `seperator`, which acts as the default separator for each command-line argument.
* This global seperator can be overridden by each individual command-line argument, providing flexibility and control.
* Command-line arguments themselves are detailed under `items`, allowing for fine-grained configuration.

For a deeper understanding of seperator, consult the [seperator section](#specifying-seperator). Likewise, for comprehensive information on items, refer to the [items section](#specifying-items).

Prepare yourself, brave adventurer, for within the realm of arguments lies the key to unlocking limitless possibilities in command-line configuration. 🚀🔧


```yaml
arguments:
  seperator:
    arguments: "?"
    argumentvalue: "="
  items:
  - name: "Server Base Argument"
    type: "base"
    mandatory: true
    argument: "server"
    seperator:
      arguments: " "
  - name: "Map"
    type: "enum"
    mandatory: true
    argument: ""
    seperator:
      argumentvalue: ""
    items:
      - name: "Arena [FFA]"
        value: "AOCFFA-Arena3_p"
      - ...
```

### Specifying `seperator`

Diving into the Depths of seperator a fundamental component of command-line configuration. In this realm, we define how command-line arguments are parsed, setting the stage for seamless execution of our commands.

But what exactly does this entail? Let's uncover the secrets:

* A seperator dictates how command-line arguments are separated and parsed.
* It defines the delimiter between each command-line argument and specifies the delimiter between the argument and its value, if applicable.
* The global seperator serves as the default for all command-line arguments, ensuring consistency in parsing.
* If no global seperator is provided, the Space-Seperator 👾 will be utilized by default.

```yaml
                                  👾 Space-Seperator
seperator:                        seperator:
  arguments: "?"                    arguments: " "
  argumentvalue: "="                argumentvalue: " "
```


#### Individual Override per Item

Each command-line argument has the power to override the global seperator, granting autonomy in parsing. Take control by specifying the seperator for each item, as demonstrated below:

```yaml
arguments:
  seperator:
    arguments: "?"
    argumentvalue: "="
  items:
  - type: "base"
    argument: "server"
    seperator:
      arguments: "+"
```

### Specifying `items`

Let's unveil the mysteries of items. Items define a list of command-line arguments, each serving a specific purpose and contributing to the overall functionality of the executable.

But what secrets lie within items? Let's uncover the details:

Various predefined types are at your disposal, each tailored to specific use cases. All types share common required and optional fields, ensuring consistency in configuration:

* Required
    * `type` - Indicates the type of command-line argument.
    * `argument` - Defines the command-line argument itself.

* Optional
    * `name` - Provides a human-friendly name for the command-line argument.
    * `mandatory` - Determines if this command-line argument is compulsory.
    * `disabled` - Indicates whether this command-line argument is active during parsing.
    * `seperator` - Sets the specific seperator to apply to this command-line argument, thereby overriding the global seperator.

With items at your disposal, you are poised to embark on a journey of command-line mastery, shaping the behavior of your executable with precision and expertise. 🚀🔧

#### type: `base`

We encounter command-line directives of type base as stand alone, devoid of any accompanying values. A base argument serves as a fundamental building block, allowing for streamlined execution of specific commands. Particularly useful is the scenario where a server can be launched using the same executable as the game client, contingent upon the presence of a command-line directive such as server.

```yaml
name: "Name to display"
type: "base"
mandatory: true
disabled: false
argument: "server"
seperator:
  arguments: "?"
```

#### type: `string`

Command-line arguments of type string expand upon the foundation laid by base arguments. Distinguishing itself, a string argument boasts an essential addition: a mandatory attribute named default. This attribute serves as a pivotal tool, allowing for the specification of a `default` value.

```yaml
name: "String Argument"
type: "string"
mandatory: false
disabled: true
argument: "+set sv_hostname"
seperator:
  arguments: "?"
  argumentvalue: "="
default: "lanserver"
```

#### type: `integer`

Command-line arguments of type integer build upon the foundation established by string arguments. Setting itself apart, an integer argument introduces essential enhancements, including two mandatory fields: `minvalue` and `maxvalue`. These fields play a crucial role in validating the range of possible values assigned to the command-line argument.

```yaml
name: "Integer Argument"
type: "integer"
mandatory: false
disabled: true
argument: "+set sv_clients"
seperator:
  arguments: "?"
  argumentvalue: "="
default: 32
minvalue: 0
maxvalue: 64
```

#### type: `float`

Command-line arguments of type float build upon the foundation established by string arguments. Setting itself apart, an float argument introduces essential enhancements, including two mandatory fields: `minvalue` and `maxvalue`. These fields play a crucial role in validating the range of possible values assigned to the command-line argument.

```yaml
name: "Float Argument"
type: "float"
mandatory: false
disabled: true
argument: "+set g_gravitiy"
seperator:
  arguments: "?"
  argumentvalue: "="
default: 15.0
minvalue: 1.0
maxvalue: 99.99
```

#### type: `boolean`

Command-line arguments of type boolean extend beyond the capabilities of string arguments. Distinguishing itself, a boolean argument introduces an optional attribute: `values`. This attribute enhances flexibility by allowing the specification of potential values associated with the command-line argument.

```yaml
name: "Boolean Argument"
type: "boolean"
mandatory: false
disabled: true
argument: "+set dedicated"
seperator:
  arguments: "?"
  argumentvalue: "="
default: true
values:
  type: custom
  true: TRue
  false: FaLse
```

##### `values`

The values field determines the specific values assigned to `true` and `false` for a boolean argument. It consists of the following attributes:

* Required
    * [`type`] - Determines the data type to be utilized for the boolean values, which may be selected from the following options: __bool, boolupper, boolupperfull, integer, or custom__.

* Optional
    * [`true`] - Designates the value attributed to true, contingent upon the `type` being specified as __custom__.
    * [`false`] - Identifies the value assigned to false, applicable only when the `type` is set to __custom__.

|type|true|false|
|---|---|---|
| bool | true | false |
| boolupper | True | False |
| boolupperfull | TRUE | FALSE |
| integer | 1 | 0 |
| custom | - | - |

```yaml
values:
  type: custom
  true: "TRue"
  false: "FaLse"
```

#### type: `enum`

Command-line arguments of type `enum` extend beyond the capabilities of base arguments. A enum argument introduces an essential addition: a mandatory attribute known as `items`. This attribute houses a collection of key-value pairs, each representing a specific option. This functionality proves invaluable for scenarios where arguments must be restricted to predefined values, such as selecting maps.

Moreover, the optional `name` field provides an opportunity for customization, allowing for user-friendly labels to be displayed in the user interface.

```yaml
name: "Enum Argument"
type: "enum"
mandatory: true
disabled: false
argument: "+map"
seperator:
  arguments: "?"
  argumentvalue: "="
items:
  - name: "Ambush"
    value: "mp_convoy"
  - name: "Backlot"
    value: "mp_backlot"
  - name: "Bloc"

```

#### type: `connect`

Command-line arguments of type connect dictate how a game establishes a direct connection to a game server via command-line, but exclusively for use within the `client` field. The primary objective of a connect argument is therefore to specify the connection mechanism between the game client and the game server. To achieve this, the character `?` functions as a wildcard placeholder for the game server's IP address. 

```yaml
type: "connect"
argument: "+connect ?"
```

## Configuration ⚙️

Important configuration can be set via the [`settings.yaml`](settings.yaml) file.

| Config | Description |
|---|---|
| port | Port of the http server to listen on |
| game-config-directory | Directory of the game configuration files |
| game-file-directory | Directory of the game files |
| game-icon-directory | Directory of the game icon files |

## ⌨️ Command-line ⌨️

| Argument | Type | Default | Example | Description |
|---|---|---|---|---|
| loglevel | string | info | `lanty-server --loglevel debug` | Sets the log level (disable, trace, debug, info, warning, error, panic, fatal) |
| logenablefile | boolean | - | `lanty-server --logenablefile` | Enables logging to file |
| logfile | string | lanty.log | `lanty-server --logfile lanty.log` | Sets the log filename |
| logdir | string | log | `lanty-server --logdir log` | Sets the log directory |
| logbackups | integer | 0 | `lanty-server --logbackups 10` | Sets the number of old logs to remain |
| logfilesize | integer | 10 | `lanty-server --logfilesize 25` | Sets the size of the logs before rotating to new file |
| logage | integer | 0 | `lanty-server --logage 1` | Sets the maximum number of days to retain old logs |
| port | integer | 8090 | `lanty-server --port 8091` | Port of the http server to listen on |
| graceful-shutdown | integer | 10 | `lanty-server --graceful-shutdown 25` | Timeout in seconds to wait for a graceful shutdown of the server |
| game-config-dir | string | game-config | `lanty-server --game-config-dir config` | Directory of the game configuration files |
| game-file-dir | string | game-data | `lanty-server --game-file-dir data` | Directory of the game files |

## 🛠️ Building from Source 🥴

Before diving into the installation process, let's first tackle the task of building from source. Here's your roadmap:

* Clone the Repository: Begin by cloning this repository to your local environment.
* Run the Commands: Once cloned, execute the following commands in your terminal:

```bash
$ go build -o ./build/lanty-server.exe -v ./cmd
$ cp ./settings.yaml ./build
$ mkdir ./build/game-config ./build/game-icon ./build/game-data ./build/download
```

* Proceed to Installation: With the build process complete, you're now ready to move on to the [`installation section`](#🚀-ready-to-get-lanty-up-and-running-lets-do-this-🚀) and continue setting up Lanty. 🦖

#### ⚔️ Cross-compiling on Linux for Windows

To cross-compile on Linux for Windows, all you need is a cross-compiler. Here's how to get started:

* Set Up the Compiler: Begin by setting up your cross-compiler. We recommend using [MinGW-w64](https://www.mingw-w64.org/), but you're free to use any cross-compiler of your choice.
* Run the Commands: Once your cross-compiler is set up, execute the following commands in your terminal:

##### 🚨Attentione: If you opt to use a cross-compiler other than MinGW-w64, you'll need to adjust the `CC` variable to align with your compiler.🚨

```bash
$ GOOS=windows CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc-win32 go build -o ./build/lanty-server.exe -v ./cmd
$ cp ./settings.yaml ./build
$ mkdir ./build/game-config ./build/game-icon ./build/game-data ./build/download
```
