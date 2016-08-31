<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;

class CreateUsersTable extends Migration {

    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('users', function(Blueprint $table) {
            $table->increments('id');

            // Google Auth related
            $table->string('google_id', 250)->default('');
            $table->string('email', 250)->default('');
            $table->string('access_token', 250)->default('');

            $table->string('name');
            $table->string('username')->unique();
            $table->string('password');
            $table->boolean('admin')->nullable();
            $table->softDeletes();
            $table->timestamps();
            $table->string('remember_token')->nullable();
        });
    }


    /*
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::drop('users');
    }

}
