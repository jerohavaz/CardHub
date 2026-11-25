package cardwise.designsystem.theme

import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.runtime.CompositionLocalProvider
import androidx.compose.runtime.ReadOnlyComposable
import androidx.compose.runtime.staticCompositionLocalOf

private val LocalColors = staticCompositionLocalOf { LightColors }
private val LocalTypography = staticCompositionLocalOf { DefaultTypography }
private val LocalShapes = staticCompositionLocalOf { CardwiseShapes() }

/**
 * Entry point of the Cardwise design system.
 * Our components read from [CardwiseTheme]. MaterialTheme is only set so that
 * Material components (e.g. ModalBottomSheet) match our colors.
 */
@Composable
fun CardwiseTheme(content: @Composable () -> Unit) {
    val colors = LightColors
    CompositionLocalProvider(
        LocalColors provides colors,
        LocalTypography provides DefaultTypography,
        LocalShapes provides CardwiseShapes(),
    ) {
        MaterialTheme(
            colorScheme = lightColorScheme(
                primary = colors.blue,
                onPrimary = colors.surface,
                primaryContainer = colors.blueSoft,
                onPrimaryContainer = colors.blue,
                background = colors.screen,
                onBackground = colors.text,
                surface = colors.surface,
                onSurface = colors.text,
                surfaceVariant = colors.surface2,
                onSurfaceVariant = colors.muted,
                outline = colors.line,
                outlineVariant = colors.line,
                error = colors.danger,
            ),
            content = content,
        )
    }
}

object CardwiseTheme {
    val colors: CardwiseColors
        @Composable @ReadOnlyComposable get() = LocalColors.current
    val typography: CardwiseTypography
        @Composable @ReadOnlyComposable get() = LocalTypography.current
    val shapes: CardwiseShapes
        @Composable @ReadOnlyComposable get() = LocalShapes.current
}
