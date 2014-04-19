<?php

use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateLeasesTable extends Migration {

	/**
	 * Run the migrations.
	 *
	 * @return void
	 */
	public function up()
	{
		Schema::create('leases', function(Blueprint $table) {
            $table->increments('id');
            $table->integer('user_id')->unsigned();
            $table->string('group_id');
            $table->string('lease_ip');
            $table->string('protocol');
            $table->string('port_from');
            $table->string('port_to');
            $table->integer('expiry')->unsigned();
            $table->softDeletes();
            $table->timestamps();
            $table->foreign('user_id')->references('id')->on('users');
        });
	}

	/**
	 * Reverse the migrations.
	 *
	 * @return void
	 */
	public function down()
	{
		Schema::table('leases', function(Blueprint $table)
		{
			$table->dropForeign('leases_user_id_foreign');
		});
		Schema::drop('leases');
	}

}
