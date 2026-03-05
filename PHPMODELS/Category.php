<?php

/**
 * Created by Reliese Model.
 */

namespace App\Models;

use Illuminate\Database\Eloquent\Collection;
use Illuminate\Database\Eloquent\Model;

/**
 * Class Category
 * 
 * @property int $id
 * @property string $name
 * @property string|null $image_src
 * @property string|null $code
 * @property string|null $category_id
 * @property float|null $range_1
 * @property float|null $range_2
 * @property float|null $range_3
 * @property float|null $free_shipping
 * 
 * @property Collection|Product[] $products
 * @property Collection|Property[] $properties
 *
 * @package App\Models
 */
class Category extends Model
{
	protected $table = 'categories';
	public $timestamps = false;

	protected $casts = [
		'range_1' => 'float',
		'range_2' => 'float',
		'range_3' => 'float',
		'free_shipping' => 'float'
	];

	protected $fillable = [
		'name',
		'image_src',
		'code',
		'category_id',
		'range_1',
		'range_2',
		'range_3',
		'free_shipping'
	];

	public function products()
	{
		return $this->hasMany(Product::class);
	}

	public function properties()
	{
		return $this->hasMany(Property::class);
	}
}
