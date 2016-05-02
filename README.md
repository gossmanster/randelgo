randelgo
========

Random Mandelbrot Explorer written in Golang

docker run -p 80:80 gossmanster/randelgo

Then point a browser at port 80 on that machine. Refresh to have it explore the set.

makeDockerImage.sh will create an optimized Docker file with just a minimum Go binary in it (no depedencies and very small)


