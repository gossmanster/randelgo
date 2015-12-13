FROM busybox:latest

# Go program that explores the Mandelbrot set
ADD randelgo /

EXPOSE 1966

CMD ["/randelgo"]