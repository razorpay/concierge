<?php 
namespace LaravelDuo;
class LaravelDuo extends Duo {

    private $_AKEY = 'SUPERSECRETAKEY';
    private $_IKEY = 'IKEYFROMDUO';
    private $_SKEY = 'SKEYFROMDUO';
    private $_HOST = 'HOSTFROMDUO';

    public function get_akey()
    {
        return $this->_AKEY;
    }

    public function get_ikey()
    {
        return $this->_IKEY;
    }

    public function get_skey()
    {
        return $this->_SKEY;
    }

    public function get_host()
    {
        return $this->_HOST;
    }

} 