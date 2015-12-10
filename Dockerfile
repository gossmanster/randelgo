FROM busybox:latest

# Go program that explores the Mandelbrot set
ADD randelgo /
CMD ["/randelgo"]