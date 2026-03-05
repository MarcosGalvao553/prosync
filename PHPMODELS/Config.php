<?php

/**
 * Created by Reliese Model.
 */

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

/**
 * Class Config
 * 
 * @property int $id
 * @property string $code
 * @property string|null $description
 * @property string|null $value
 *
 * @package App\Models
 */
class Config extends Model
{
	protected $table = 'configs';
	public $timestamps = false;

	protected $fillable = [
		'code',
		'description',
		'value'
	];
}
