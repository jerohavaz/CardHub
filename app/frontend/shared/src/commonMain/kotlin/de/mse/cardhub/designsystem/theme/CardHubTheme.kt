package de.mse.cardhub.designsystem.theme

import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.runtime.CompositionLocalProvider
import androidx.compose.runtime.ReadOnlyComposable
import androidx.compose.runtime.staticCompositionLocalOf

private val LocalCardHubColors = staticCompositionLocalOf { LightCardHubColors }
private val LocalCardHubTypography = staticCompositionLocalOf { DefaultCardHubTypography }
private val LocalCardHubShapes = staticCompositionLocalOf { CardHubShapes() }

/**
 * Entry point of the CardHub design system.
 * Our components read from [CardHubTheme]. MaterialTheme is only set so that
 * Material components (e.g. ModalBottomSheet) match our colors.
 */
@Composable
fun CardHubTheme(content: @Composable () -> Unit) {
    val colors = LightCardHubColors
    CompositionLocalProvider(
        LocalCardHubColors provides colors,
        LocalCardHubTypography provides DefaultCardHubTypography,
        LocalCardHubShapes provides CardHubShapes(),
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

object CardHubTheme {
    val colors: CardHubColors
        @Composable @ReadOnlyComposable get() = LocalCardHubColors.current
    val typography: CardHubTypography
        @Composable @ReadOnlyComposable get() = LocalCardHubTypography.current
    val shapes: CardHubShapes
        @Composable @ReadOnlyComposable get() = LocalCardHubShapes.current
}
