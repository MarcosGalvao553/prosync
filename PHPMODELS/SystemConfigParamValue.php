<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class SystemConfigParamValue extends Model
{
    use HasFactory;

    protected $fillable = [
		'name',
		'code',
		'value',
        'user_id',
        'system_config_param_id',
        'is_config'
	];


    public function system_config_param()
	{
		return $this->belongsTo(SystemConfigParam::class);
	}

    public function user()
	{
		return $this->belongsTo(User::class);
	}

	public function get($key)
    {
        $key = $this->where('name', $key)->first();
        return $key->value;
    }
}
