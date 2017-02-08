<?php

namespace App\Console\Commands;

use App\Http\Controllers\HomeController;
use Illuminate\Console\Command;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputOption;

class LeaseManager extends Command
{
    /**
     * The console command name.
     *
     * @var string
     */
    protected $name = 'custom:leasemanager';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Manages termination & removal of expired/obsolete leases.';

    /**
     * Execute the console command.
     * Creates a HomeController and calls its cleanLeases function which terminates all expired leases.
     *
     * @return mixed
     */
    public function handle()
    {
        $home_controller = new HomeController();
        $cleaner = $home_controller->cleanLeases();
        if ($cleaner) {
            $this->info($cleaner);
        }
    }

    /*
     * Get the console command arguments.
     *
     * @return array
     */
    /*protected function getArguments()
    {
        return array(
            array('example', InputArgument::REQUIRED, 'An example argument.'),
        );
    }*/

    /*
     * Get the console command options.
     *
     * @return array
     */
    /*protected function getOptions()
    {
        return array(
            array('example', null, InputOption::VALUE_OPTIONAL, 'An example option.', null),
        );
    }*/
}
