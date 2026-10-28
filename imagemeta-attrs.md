# imagemeta v1.1.0 — available attributes

`imagemeta.Decode(r)` returns `exif.Exif`. Go field paths are shown below. ✅ = already saved in `service.Metadata`.

## Top level (`ex.`)
| Field | Type | Note |
|---|---|---|
| `ImageType` ✅ | imagetype.ImageType | jpeg, heic, png, cr2, cr3, nef, arw, dng… |
| `CameraSerial` | string | from the maker note |
| `CameraMakeID` | makernote.CameraMake | maker as an enum |
| `IFD0` | IFD0Tag | main image tags |
| `ExifIFD` | ExifIFDTags | camera / shot tags |
| `GPS` | GPSInfo | location (read through methods) |
| `IFD1`, `IFD2` | *ImageIFD | thumbnail / extra image |
| `SubIFDs` | [8]*ImageIFD | previews embedded in RAW files |
| `MakerNote` | makernote.Info | vendor data: Apple, Canon, Nikon, Sony, Panasonic, DNG |

## IFD0 (`ex.IFD0.`): main image
- **Camera/author:** `Make` ✅, `Model` ✅, `Software`, `Artist`, `Copyright`, `ImageDescription`, `Rating`
- **Dimensions:** `ImageWidth`, `ImageHeight`, `XResolution`, `YResolution`, `ResolutionUnit`, `Orientation`
- **Encoding:** `Compression`, `PhotometricInterpretation`, `SamplesPerPixel`, `BitsPerSample`, `PlanarConfiguration`, `YCbCrPositioning`, `SubfileType`
- **Offsets:** `ImageOffset`, `ImageLength`, `ThumbnailOffset`, `ThumbnailLength`
- **Time:** `ModifyDate`, `OffsetTime`, `SubSecTime`

## ExifIFD (`ex.ExifIFD.`): shot details
- **Time:** `DateTimeOriginal` ✅, `CreateDate`, `OffsetTimeOriginal`, `OffsetTimeDigitized`, `SubSecTimeOriginal`, `SubSecTimeDigitized`
- **Exposure:** `ExposureTime`, `FNumber`, `ApertureValue`, `MaxApertureValue`, `ShutterSpeedValue`, `ISOSpeedRatings`, `RecommendedExposureIndex`, `SensitivityType`, `ExposureIndex`, `ExposureProgram`, `ExposureMode`, `ExposureBias`, `BrightnessValue`, `MeteringMode`, `Flash`
- **Lens:** `LensMake`, `LensModel`, `LensSerial`, `LensInfo` (min/max focal length and aperture), `FocalLength`, `FocalLengthIn35mmFormat`, `DigitalZoomRatio`, `SubjectDistance`, `SubjectDistanceRange`, `SubjectArea`
- **Body/owner:** `BodySerialNumber`, `CameraOwnerName`, `UserComment`
- **Dimensions:** `PixelXDimension` ✅, `PixelYDimension` ✅, `RelatedImageWidth`, `RelatedImageHeight`, `FocalPlaneXResolution`, `FocalPlaneYResolution`, `FocalPlaneResolutionUnit`
- **Look:** `ColorSpace`, `WhiteBalance`, `LightSource`, `Contrast`, `Saturation`, `Sharpness`, `GainControl`, `Gamma`, `CustomRendered`, `SceneCaptureType`, `SceneType`, `SensingMethod`, `FileSource`, `CompositeImage`
- **Format info:** `ExifVersion`, `FlashpixVersion`, `ComponentsConfiguration`, `CompressedBitsPerPixel`, `InteropIndex`, `InteropVersion`

## GPS (`ex.GPS.`): methods, since the fields are private
- **Position:** `Latitude()` ✅, `Longitude()` ✅, `Altitude()`
- **Time:** `GPSTimestamp()`
- **Accuracy:** `DOP()`, `HPositioningError()`, `Satellites()`, `Status()`, `MeasureMode()`, `Differential()`, `MapDatum()`, `ProcessingMethod()`, `AreaInformation()`
- **Motion/direction:** `SpeedWithRef()`, `TrackWithRef()`, `ImgDirectionWithRef()`
- **Destination:** `DestLatitude()`, `DestLongitude()`, `DestLatitudeSigned()`, `DestLongitudeSigned()`, `DestBearingWithRef()`, `DestDistanceWithRef()`
- **Field:** `GPSVersion`

## ImageIFD (`IFD1`, `IFD2`, `SubIFDs[i]`)
`Make`, `Model`, `Software`, `ImageDescription`, `ImageWidth`, `ImageHeight`, `XResolution`, `YResolution`, `ResolutionUnit`, `Compression`, `Orientation`, `SubfileType`, `ImageOffset`, `ImageLength`, `ModifyDate`, `OffsetTime`, `SubSecTime`

## Other functions
- `DecodeJPEG`, `DecodeHeif`, `DecodePng`, `DecodeTiff`, `DecodeCR2`, `DecodeCR3`, `DecodeARW`, `DecodeCRW`: decode a known format
- `PreviewCR3(r)`: returns the embedded JPEG preview from a CR3 file
- `exif.Exif` implements `MarshalJSON`, so `json.Marshal(ex)` returns every field above

---

# By platform

Every platform fills the standard groups (`IFD0`, `ExifIFD`, `GPS`) when the file has EXIF. What differs is the **vendor maker note** (`ex.MakerNote.<Vendor>`), which the lib parses only for the vendors below. The "usually filled" lists are typical camera behaviour, not something the lib guarantees.

## iPhone / iOS (`MakerNote.Apple`), usually HEIC or JPEG
- **Standard, usually filled:** Make/Model, LensMake/LensModel, DateTimeOriginal + offsets, exposure/ISO/aperture, FocalLengthIn35mmFormat, full GPS (incl. altitude, speed, direction, HPositioningError)
- **Apple-only:**
  - capture: `ImageCaptureType`, `CameraType`, `HDRImageType`, `HDRGain`, `HDRHeadroom`, `SemanticStyle*`, `SceneFlags`
  - focus/AE: `AFStable`, `AFConfidence`, `AFMeasuredDepth`, `FocusPosition`, `FocusDistanceRange`, `AEStable`, `AETarget`, `AEAverage`, `OISMode`
  - IDs: `ContentIdentifier` (links a Live Photo to its video), `LivePhotoVideoIndex`, `BurstUUID`, `ImageUniqueID`, `PhotoIdentifier`, `ImageCaptureRequestID`
  - sensor: `AccelerationVector`, `ColorTemperature`, `SignalToNoiseRatio`, `LuminanceNoiseAmplitude`, `RunTime`

## Android (Samsung, Pixel, Xiaomi…), usually JPEG or HEIC
- **No maker-note parser**, so only standard fields are available
- **Usually filled:** Make/Model, Software, DateTimeOriginal + offsets, exposure/ISO/aperture/FocalLength, Flash, WhiteBalance, GPS (lat/lng/alt/timestamp)
- **Often empty:** LensMake/LensModel, BodySerialNumber

## Canon (`MakerNote.Canon`), CR2 / CR3 / CRW / JPEG
- **Body/IDs:** `FirmwareVersion`, `OwnerName`, `SerialNumber`, `InternalSerialNumber`, `ModelID`, `FileNumber`, `ImageUniqueID`, `BatteryType`
- **Lens:** `LensModel`, `LensInfo`, `CanonFocalLength`, `AFMicroAdj`
- **Shot:** `CanonCameraSettings`, `CanonShotInfo`, `ColorTemperature`, `ProcessingInfo`, `CustomPictureStyleFileName`, `HDRInfo`, `MultiExp`
- **Focus/sensor:** `AFInfo`, `AFConfig`, `FaceDetect`, `SensorInfo`, `AspectInfo`, `TimeInfo`
- **Also:** `ex.CameraSerial`; `imagemeta.PreviewCR3()` returns the embedded JPEG

## Nikon (`MakerNote.Nikon`), NEF / JPEG
- **Body:** `SerialNumber`, `ShutterCount`, `MechanicalShutterCount`, `ImageCount`, `DeletedImageCount`, `PowerUpTime`, `WorldTime`
- **Lens:** `Lens`, `LensType`, `LensFStops`, `VRInfo`, `ManualFocusDistance`
- **Shot:** `ISO`, `ISOSetting`, `ISOInfo`, `Quality`, `ColorMode`, `WhiteBalance`, `Sharpness`, `ActiveDLighting`, `VignetteControl`, `ShootingMode`, `ShutterMode`, `SilentPhotography`, `FlashSetting/Type/Mode`
- **Focus:** `FocusMode`, `AFInfo`, `AFInfo2`, `AFTune`
- **RAW preview:** `ex.SubIFDs` (JpgFromRaw)

## Sony (`MakerNote.Sony`), ARW / JPEG
- **Body/lens:** `SerialNumber`, `SonyModelID`, `LensType`, `LensSpec`, `Teleconverter`
- **Shot:** `CreativeStyle`, `SceneMode`, `Quality`, `JPEGQuality`, `ReleaseMode`, `SequenceNumber`, `HDR`, `DynamicRangeOptimizer`, `ImageStabilization`, `Rating`
- **Tone:** `Contrast`, `Saturation`, `Sharpness`, `Shadows`, `Highlights`, `Fade`, `Clarity`
- **White balance/colour:** `WhiteBalance`, `WhiteBalanceFineTune`, `WBShiftABGM`, `ColorTemperature`
- **Focus:** `FocusMode`, `AFAreaMode`, `AFTracking`, `FocusLocation`, `AFInfo`
- **Noise/lens corrections:** `HighISONoiseReduction`, `MultiFrameNoiseReduction`, `VignettingCorrection`, `DistortionCorrectionSetting`
- **Sub-blocks:** `CameraInfo*`, `CameraSettings*`, `ShotInfo`, `Tag94xx`

## Panasonic / Lumix (`MakerNote.Panasonic`), RW2 / JPEG
- **Body:** `InternalSerialNumber`, `FirmwareVersion`, `BatteryLevel`, `TimeSincePowerOn`, `TravelDay`
- **Shot:** `ImageQuality`, `ShootingMode`, `SceneMode`, `BurstMode`, `SequenceNumber`, `SelfTimer`, `MacroMode`, `ImageStabilization`, `ColorEffect`, `ContrastMode`, `NoiseReduction`
- **White balance:** `WhiteBalance`, `WhiteBalanceBias`, `WBRed/Green/BlueLevel`
- **Focus:** `FocusMode`, `AFAreaMode`, `AFPointPosition`

## DNG / Adobe (`MakerNote.DNG`)
- `DNGVersion`, `DNGBackwardVersion`, `CameraModel`, `OriginalRawFileName`, `ProfileName`, `BestQualityScale`, `AdobeData`

## Web / edited / screenshots (PNG, re-saved JPEG, social-media downloads)
- **Usually no EXIF at all:** `Decode` fails, and `service.ExtractMetadata` falls back to `image.DecodeConfig` (format + width/height only)
- **PNG:** only the `eXIf` chunk, if one is present
- **Edited files:** often only `IFD0.Software` plus `ModifyDate`
