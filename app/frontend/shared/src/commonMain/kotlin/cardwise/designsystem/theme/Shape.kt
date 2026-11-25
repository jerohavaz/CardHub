package cardwise.designsystem.theme

import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.runtime.Immutable
import androidx.compose.ui.unit.dp

/** Corner radii from the Figma template. */
@Immutable
data class CardwiseShapes(
    /** Small badges in tiles, stepper buttons, plus button */
    val xs: RoundedCornerShape = RoundedCornerShape(7.dp),
    /** Inner tabs in the segmented control */
    val sm: RoundedCornerShape = RoundedCornerShape(10.dp),
    /** Stepper, segmented control */
    val md: RoundedCornerShape = RoundedCornerShape(12.dp),
    /** Buttons, search field, info boxes, icon button (square) */
    val lg: RoundedCornerShape = RoundedCornerShape(14.dp),
    /** List rows, stats */
    val xl: RoundedCornerShape = RoundedCornerShape(16.dp),
    /** Card tiles, panels */
    val card: RoundedCornerShape = RoundedCornerShape(19.dp),
    /** Hero banner */
    val hero: RoundedCornerShape = RoundedCornerShape(22.dp),
    /** Floating navigation, bottom sheets */
    val sheet: RoundedCornerShape = RoundedCornerShape(28.dp),
    /** Chips, badges */
    val pill: RoundedCornerShape = RoundedCornerShape(50),
)

/** Spacing – the design consistently uses a 20 dp horizontal margin. */
object Spacing {
    val screenHorizontal = 20.dp
    val xs = 4.dp
    val sm = 8.dp
    val md = 12.dp
    val lg = 16.dp
    val xl = 20.dp
    val xxl = 24.dp
}
