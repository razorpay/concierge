<?php

use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateInvitesTable extends Migration {

	/**
	 * Run the migrations.
	 *
	 * @return void
	 */
	public function up()
	{
		Schema::create('invites', function(Blueprint $table) {
            $table->increments('id');
            $table->integer('user_id')->unsigned();
            $table->string('group_id');
            $table->string('protocol');
            $table->string('port_from');
            $table->string('port_to');
            $table->integer('expiry')->unsigned();
            $table->string('email')->nullable();
            $table->string('token')->unique();
            $table->softDeletes();
            $table->timestamps();
            $table->foreign('user_id')->references('id')->on('users')->onDelete('cascade');;
        });
	}

	/**
	 * Reverse the migrations.
	 *
	 * @return void
	 */
	public function down()
	{
		Schema::table('invites', function(Blueprint $table)
		{
			$table->dropForeign('invites_user_id_foreign');
		});
		Schema::drop('invites');
	}

}
