# ghash
Hash function with image output (identicon)

This command line tools create an image based on its input.

The goal is to create "something" easier to verify and remember than hexa string of typical crytpographic hash.

To do that ghash use sha512 3 times, to generate 22 numbers between 0
and 1.0, to create an image using 2 julia fractal curves mixed
together and a 4 colors gradient.

Option "-9" will create a new image using sha3 and bcrypt as image parameter. This enable the use of the 2 images as a strong visual hash.

## example

$ echo "Hello world!" | ghash -o helloworld.png

![Doc](https://github.com/nicolasboulay/ghash/raw/master/example/helloworld.png "helloworld.png")

$ cat ./gmic | ./ghash -o gmic.jpg 

![Doc](https://github.com/nicolasboulay/ghash/raw/master/example/gmic.jpg "gmic.jpg")

$ echo "one small icon" | ./ghash -size 32 -o icon32.png

![Doc](https://github.com/nicolasboulay/ghash/raw/master/example/icon32.png "icon32.png")

$ echo large | ./ghash -size 500 -o ../example/large.jpg

![Doc](https://github.com/nicolasboulay/ghash/raw/master/example/large.jpg "large.jpg")

$ echo "Hello world!" | ghash -9 -o helloworld.png -o2 helloworld2.png

![Doc](https://github.com/nicolasboulay/ghash/raw/master/example/helloworld.png "helloworld.png")
![Doc](https://github.com/nicolasboulay/ghash/raw/master/example/helloworld2.png "helloworld2.png")


## More information

Humain brain are more sensitive to shape than color: the image generator must not use too much color. But the image generator should generate enough different images to avoid collision (you add some characters to a texte, to produce the same image than an other texte). Ghash can generate around 100^22 images, but only 10^22 could be easly differentiated by the naked eyes.  

Attack could be done by brute force using the generator, using a
metric like the psnr, to find 2 similar images under a threasold, but
the generator take 200ms per image.

Ghash use "go" langage to generate the SHA512 internal hash, and a gmic script to generate the image ( http://gmic.eu ). The gmic script could be rewritten in C++ using the CImg.h header, for faster generation (around 10x). The actual script take less than a second. 

The "-9" option rehash the input with the sha3-512 algorithme injected inside a bcrypt hash, then the result is hashed 3 time with sha3 to create the 22 parameters. 
