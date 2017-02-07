<?php

namespace App\Exception;

use Exception;
use Response;
use Illuminate\Database\Eloquent\ModelNotFoundException;
use Symfony\Component\HttpKernel\Exception\HttpException;
use Symfony\Component\HttpKernel\Exception\NotFoundHttpException;
use Illuminate\Session\TokenMismatchException;
use Illuminate\Contracts\Encryption\DecryptException;
use Illuminate\Foundation\Exceptions\Handler as ExceptionHandler;

class Handler extends ExceptionHandler
{
    const SERVER_ERROR       = 'Internal Server Error';

    const METHOD_NOT_ALLOWED = 'Method not allowed';
    const ENTITY_NOT_FOUND   = 'Not Found';

    const RESPONSE_404 = [
        'success' => false,
        'errors' => [
            self::ENTITY_NOT_FOUND
        ]
    ];

    /**
     * A list of the exception types that should not be reported.
     *
     * @var array
     */
    protected $dontReport = [
        DecryptException::class,
        HttpException::class,
        ModelNotFoundException::class,
        NotFoundHttpException::class,
        TokenMismatchException::class,
    ];

    protected function isCritical(Exception $e)
    {
        foreach ($this->dontReport as $type)
        {
            if ($e instanceof $type)
            {
                return false;
            }
        }

        return true;
    }
}
