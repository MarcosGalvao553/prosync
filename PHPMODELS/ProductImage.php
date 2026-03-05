<?php

/**
 * Created by Reliese Model.
 */

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

/**
 * Class ProductImage
 * 
 * @property int $id
 * @property int $image_type
 * @property string $image_src
 * @property int $product_id
 * @property string|null $Image_src_small
 * 
 * @property Product $product
 *
 * @package App\Models
 */
class ProductImage extends Model
{
	protected $table = 'product_images';
	public $timestamps = false;

	protected $casts = [
		'image_type' => 'int',
		'product_id' => 'int'
	];

	protected $fillable = [
		'image_type',
		'image_src',
		'product_id',
		'Image_src_small'
	];

	public function product()
	{
		return $this->belongsTo(Product::class);
	}
}
