<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class SystemConfigParam extends Model
{
    use HasFactory;

    protected $fillable = [
		'name',
		'description',
        'show_to_user',
        'code',
        'system_config_id'
	];

	public function system_config_params_values()
	{
		return $this->hasMany(SystemConfigParamValue::class);
	}

    public function system_config()
	{
		return $this->belongsTo(SystemConfig::class);
	}
}
