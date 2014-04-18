<?php

class BaseController extends Controller {

	/**
	 * Setup the layout used by the controller.
	 *
	 * @return void
	 */
	protected function setupLayout()
	{
		if ( ! is_null($this->layout))
		{
			$this->layout = View::make($this->layout);
		}
	}

	/**
	 * Get the list of all security groups
	 */
	public function getIndex()
	{

		$ec2 = App::make('aws')->get('ec2');
		$security_groups=$ec2->describeSecurityGroups(array(
			'Filters' => array(
				array(
					'Name' => 'vpc-id',
					'Values' => array('vpc-4ff9012a'),
                ),
            ),
        ));

		$security_groups=$security_groups['SecurityGroups'];

        return View::make('getIndex')->with('security_groups', $security_groups);
	}

}
