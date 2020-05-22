# Brazilian postal codes API

An API for Brazilian postal codes, where will be available to consult the addresses by zipcode.

## API Endpoint

* GET `/zipcode/<zipcode>` Get the correspondence address based on zip code.  

```shell
> curl localhost:8000/zipcode/85801000

{
    "federative_unit": "PR",
    "city": "Cascavel",
    "neighborhood": "Centro",
    "address_name": "Avenida Brasil",
    "complement": "- de 5623 a 6869 - lado ímpar",
    "zipcode": "85801000",
    "created_at": "2020-05-22T11:56:51.735545Z",
    "updated_at": "2020-05-22T11:56:51.7355607Z"
}
```

## License

The MIT License (MIT)
