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
        $this->_AKEY = config('duo.akey');
        $this->_IKEY = config('duo.ikey');
        $this->_SKEY = config('duo.skey');
        $this->_HOST = config('duo.host');
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
