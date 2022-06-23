package main

import "fmt"

var san_cnf = `[req]
default_bits = 2048
distinguished_name = req_distinguished_name
req_extensions = req_ext
x509_extensions = v3_req
prompt = no
[req_distinguished_name]
countryName = %s
stateOrProvinceName = %s
localityName = %s
organizationName = %s
OU=%s
commonName = %s
[req_ext]
subjectAltName = @alt_names
[v3_req]
subjectAltName = @alt_names
[alt_names]
IP.1 = %s
`

func get_san_cnf() string {
	var cn, sn, ln, on, ou string

	InfoLogger.Println("Please enter cert Info:")
	fmt.Print("Please enter countryName (2 letter):")
	_, err := fmt.Scan(&cn)
	handle_error(err)
	fmt.Print("Please enter stateOrProvinceName:")
	_, err = fmt.Scan(&sn)
	handle_error(err)
	fmt.Print("Please enter localityName:")
	_, err = fmt.Scan(&ln)
	handle_error(err)
	fmt.Print("Please enter Organization Name:")
	_, err = fmt.Scan(&on)
	handle_error(err)
	fmt.Print("Please enter Organizational Unit:")
	_, err = fmt.Scan(&ou)
	handle_error(err)

	return fmt.Sprintf(san_cnf, cn, sn, ln, on, ou, get_ip(), get_ip())
}
