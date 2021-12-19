# **M1MC**

Apple have used arm64 architecture on their new Macs.
But Minecraft have not support them officially.
So I make this program to make Minecraft run natively with [HMCL](https://github.com/huanghongxun/HMCL)

# **How to use**
1. Go to [releases page](https://github.com/Dreamail/M1MC/releases), download M1MC.zip under "Assets".
2. Unzip M1MC.zip, open M1MC folder, right click "m1mc" with Alt key, click Copy "m1mc" as Pathname.
3. Open HMCL and click Settings, paste into "Wrapper command" blank.
4. Custom "Native Library Path" with M1MC/natives.
5. Done! Enjoy your play!

# How to Build
1. Clone this repository.
2. run command `go build -o m1mc main.go`.
3. Download lwjgl macos arm64 library and unzip necessary jars into lwjgl floder, native library into natives.