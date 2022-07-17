![brcepgopher](docs/img/gopher.png)

# brcep 

API for accessing information from Brazilian CEPs. The central idea is not to be dependent on a specific API, but to have the ease of accessing __brcep__ and it is in charge of consulting various sources and returning the CEP information quickly and easily.

Currently we support API queries to [ViaCEP](http://viacep.com.br), [CEPAberto](http://cepaberto.com) and [Correios](https://apps.correios.com.br/). Your help is welcome to implement the `CepApi` interface and introduce new APIs support.

![brcep](docs/img/brcep.png)

### Package

I forked this repository to transform it in a package that to import directly into my golang projects, instead of using the original sidecar pattern. 

Some other changes:
* I removed the prefered API configuration. Instead, this library returns the result that comes first. (I just want it to be faster)
* Removed cache and saved to dict (for most applications, ceps do not change, but could add this back latter on)
* Removed most dependencies (i.e. logger and server, they have no use)
* Renamed package handle to cep ( Applications should see this module as a cep) 
* and api to basecep

```


