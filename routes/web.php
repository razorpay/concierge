<?php

/*
|--------------------------------------------------------------------------
| Web Routes
|--------------------------------------------------------------------------
|
| Here is where you can register web routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| contains the "web" middleware group. Now create something great!
|
*/

// Route::get('/', function () {
//     return view('welcome');
// });

Route::group(['middleware' => ['web']], function () {
    //Routes for logged in user
    Route::group(['middleware' => ['auth']], function () {
        Route::get('/groups', 'GroupController@getGroups');

        Route::get('/manage/{group_id}', 'GroupController@getManage');

        Route::post('/manage/{group_id}', ['before'=>'csrf', 'uses'=>'GroupController@postManage']);

        Route::post('/manage/{group_id}/renew', ['before'=>'csrf', 'uses'=>'GroupController@postRenew']);

        Route::post('/manage/{group_id}/terminate', ['before'=>'csrf', 'uses'=>'GroupController@postTerminate']);

        Route::get('/logout', 'UserController@getLogout');

        Route::get('/password', 'UserController@getPassword');

        Route::post('/password', ['before'=>'csrf', 'uses'=>'UserController@postPassword']);

        //Routes for site admin
        Route::group(['before' => ['admin']], function () {
            Route::get('/users', 'UserController@getUsers');

            Route::post('/users', ['before'=>'csrf', 'uses'=>'UserController@postUsers']);

            Route::get('/users/add', 'UserController@getAddUser');

            Route::post('/users/add', ['before'=>'csrf', 'uses'=> 'UserController@postAddUser']);

            Route::get('/user/{id}/edit', 'UserController@getEditUser');

            Route::post('/user/{id}/edit', 'UserController@postEditUser');
        });
    });

    //Routes for non-logged in user
    Route::group(['middleware' => 'guest'], function () {
        Route::get('/', 'UserController@getIndex')->name('/');

        Route::get('/status', 'UserController@getStatus');
    });

    Route::get('/invite/{token}', 'UserController@getInvite');
});

// Non-browser routes
Route::group(['middleware' => ['cron']], function () {
    Route::any('/cron', 'GroupController@cleanLeases');
});
