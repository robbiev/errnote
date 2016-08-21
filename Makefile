.PHONY: clean run

errnote: clean errnote.c
	cc `pkg-config --cflags gtk+-3.0` errnote.c -o errnote `pkg-config --libs gtk+-3.0`

clean:
	rm -f errnote

run: errnote
	./errnote
