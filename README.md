# M1MC

Apple is using arm64 architecture on their new Macs.

But Minecraft hasn't supported macOS arm64 officially yet.

So I made this program to make Minecraft run natively with [HMCL](https://github.com/huanghongxun/HMCL)

## How to use

1. Go to [releases page](https://github.com/Dreamail/M1MC/releases), download `M1MC.zip` under "Assets".
2. Unzip M1MC.zip, open M1MC folder, right click "m1mc" with Alt key, click Copy "m1mc" as Pathname.
3. Open HMCL and click Settings, paste into "Wrapper command" field.
4. Custom "Native Library Path" with any path, because M1MC will modify it when game start and HMCL won't work if you don't use custom native lib path.
5. Done! Enjoy your play!

## How to Build

1. Clone this repository.
2. Run command `go build -o m1mc main.go`.
3. ~~Download [LWJGL macOS arm64](https://www.lwjgl.org/customize) and [`java-objc-bridge`](https://mvnrepository.com/artifact/ca.weblite/java-objc-bridge) and unzip necessary jars into `libraries` folder, native libraries into `natives`.~~ No longer need! M1MC will download libraries when game start at first time. 
