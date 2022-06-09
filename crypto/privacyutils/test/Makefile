SRC_FILE=main.cpp
PROG=privacy_test

$(PROG): privacy
	g++ --std=c++11 $(SRC_FILE) -o $@ -L../ -lbuchain_secp256k1 -lssl -lcrypto -I../ 

.PHONY: privacy
privacy:
	cmake ../
	make -C ../


.PHONY: clean run
clean:
	rm -f $(PROG)

run: $(PROG)
	./$(PROG) --perf-test 100
