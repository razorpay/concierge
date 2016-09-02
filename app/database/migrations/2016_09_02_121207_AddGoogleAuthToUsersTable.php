<?php

use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddGoogleAuthToUsersTable extends Migration {

	/**
	 * Run the migrations.
	 *
	 * @return void
	 */
	public function up()
	{
        Schema::table('users', function($table)
        {
            // Google Auth related
            $table->string('google_id', 250)->default('');
            $table->string('email', 250)->default('');
            $table->string('access_token', 250)->default('');
        });
	}

	/**
	 * Reverse the migrations.
	 *
	 * @return void
	 */
	public function down()
	{
		Schema::table('users', function($table)
        {
            $table->dropColumn('google_id');
            $table->dropColumn('email');
            $table->dropColumn('access_token');
        });
	}

}
