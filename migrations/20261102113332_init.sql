-- Create "images" table
CREATE TABLE "images" (
  "id" bigserial NOT NULL,
  "key" text NOT NULL,
  "format" text NULL,
  "mime_type" text NULL,
  "width" bigint NULL,
  "height" bigint NULL,
  "size" bigint NULL,
  "make" text NULL,
  "model" text NULL,
  "camera_serial" text NULL,
  "lens_make" text NULL,
  "lens_model" text NULL,
  "software" text NULL,
  "artist" text NULL,
  "copyright" text NULL,
  "orientation" text NULL,
  "taken_at" timestamptz NULL,
  "modified_at" timestamptz NULL,
  "exposure_time" text NULL,
  "f_number" numeric NULL,
  "iso" bigint NULL,
  "focal_length" numeric NULL,
  "focal_length35mm" numeric NULL,
  "exposure_bias" text NULL,
  "flash_fired" boolean NULL,
  "lat" numeric NULL,
  "lng" numeric NULL,
  "altitude" numeric NULL,
  "exif" jsonb NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_images_key" to table: "images"
CREATE UNIQUE INDEX "idx_images_key" ON "images" ("key");
-- Create index "idx_images_taken_at" to table: "images"
CREATE INDEX "idx_images_taken_at" ON "images" ("taken_at");
