# ghash
Hash function with image output

This command line tools create an image based on its input.

The goal is to create "something" easier to verify and remember than hexa string of typical crytpographic hash.

To do that ghash use sha512 3 times, to generate 22 number between 0 and 1.0, to create an image using 2 julia fractal curves mixed together and a 3 color gradient.

## example

$ echo "Hello world!" | ghash -o helloworld.png

![Doc](https://github.com/nicolasboulay/ghash/raw/master/example/helloworld.png "helloworld.png")

## More information

Humain brain are more sensible to shape than color : the image generator must not use too much color. But the image genetor should generate enough different images, to avoid collision (you add a character to a texte, to produce the same image than an other texte). Ghash can generate around 100^22 image, but only 10^22 could be easly differentiate by the naked eyes.  

Attack could be done by brute forced the generator, using a metric like the psnr, to find 2 similar images under a threasold.

Ghash use "go" langage to generate the SHA512 internal hash, and a gmic script to generate the image ( http://gmic.eu ). The gmic part could be rewritten in C++ using the CImg.h header, for faster generation (around 10x). The actual script take less than a second. 
