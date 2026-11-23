package de.mse.cardhub.designsystem.components

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.ColumnScope
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.navigationBarsPadding
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.ModalBottomSheet
import androidx.compose.material3.Text
import androidx.compose.material3.rememberModalBottomSheetState
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import de.mse.cardhub.designsystem.icons.CardwiseIcon
import de.mse.cardhub.designsystem.theme.CardwiseTheme

/**
 * Bottom sheet with handle, title, close button and content.
 * Based on Material3 [ModalBottomSheet] so swipe and back gestures work on both platforms.
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun CardwiseBottomSheet(
    title: String,
    onDismiss: () -> Unit,
    content: @Composable ColumnScope.() -> Unit,
) {
    val colors = CardwiseTheme.colors
    ModalBottomSheet(
        onDismissRequest = onDismiss,
        sheetState = rememberModalBottomSheetState(skipPartiallyExpanded = true),
        containerColor = colors.screen,
        shape = RoundedCornerShape(topStart = 28.dp, topEnd = 28.dp),
        dragHandle = {
            Box(
                Modifier
                    .padding(top = 9.dp)
                    .size(width = 40.dp, height = 5.dp)
                    .background(colors.line, RoundedCornerShape(5.dp)),
            )
        },
    ) {
        Column(Modifier.fillMaxWidth().padding(horizontal = 21.dp).padding(bottom = 30.dp).navigationBarsPadding()) {
            Row(Modifier.fillMaxWidth().padding(top = 19.dp), verticalAlignment = Alignment.CenterVertically) {
                Text(title, style = CardwiseTheme.typography.title, color = colors.text, modifier = Modifier.weight(1f))
                CardwiseIconButton(CardwiseIcon.Close, onClick = onDismiss, contentDescription = "Schließen")
            }
            content()
        }
    }
}

/** Confirmation sheet with cancel and confirm action. */
@Composable
fun ConfirmSheet(
    title: String,
    message: String,
    confirmText: String,
    onConfirm: () -> Unit,
    onDismiss: () -> Unit,
    destructive: Boolean = true,
    cancelText: String = "Abbrechen",
) {
    CardwiseBottomSheet(title = title, onDismiss = onDismiss) {
        Text(
            message,
            style = CardwiseTheme.typography.caption,
            color = CardwiseTheme.colors.muted,
            modifier = Modifier.padding(top = 16.dp),
        )
        Row(Modifier.fillMaxWidth().padding(top = 18.dp), horizontalArrangement = Arrangement.spacedBy(9.dp)) {
            CardwiseButton(cancelText, onClick = onDismiss, variant = ButtonVariant.Neutral, modifier = Modifier.weight(1f))
            CardwiseButton(
                confirmText,
                onClick = onConfirm,
                variant = if (destructive) ButtonVariant.Danger else ButtonVariant.Primary,
                modifier = Modifier.weight(1f),
            )
        }
    }
}
