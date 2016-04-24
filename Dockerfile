FROM busybox:latest

# Go program that explores the Mandelbrot set
ADD randelgo /

EXPOSE 80

CMD ["/randelgo"]