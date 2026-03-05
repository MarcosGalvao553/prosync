<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class SystemConfig extends Model
{
    use HasFactory;

    protected $fillable = [
		'name',
		'description'
	];

	public function system_config_params()
	{
		return $this->hasMany(SystemConfigParam::class);
	}
}
