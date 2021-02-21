apis:
	cd $(dirname $0)
	git clone https://github.com/hu-tao-supremacy/api.git apis
	python3 sym.py

apisgen:
	cd $(dirname $0)
	rm -rf apis
	rm -rf hts
	git clone https://github.com/hu-tao-supremacy/api.git apis
	python3 sym.py
