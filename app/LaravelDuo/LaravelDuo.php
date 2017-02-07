<?php

namespace App\LaravelDuo;

use Duo\Web as DuoWeb;

class LaravelDuo extends DuoWeb
{
    private $_AKEY;
    private $_IKEY;
    private $_SKEY;
    private $_HOST;

    public function __construct() {
        $this->_AKEY = config('custom_config.duo_akey');
        $this->_IKEY = config('custom_config.duo_ikey');
        $this->_SKEY = config('custom_config.duo_skey');
        $this->_HOST = config('custom_config.duo_host');
    }

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
