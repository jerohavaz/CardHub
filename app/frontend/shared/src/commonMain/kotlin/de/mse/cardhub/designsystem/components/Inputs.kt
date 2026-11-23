package de.mse.cardhub.designsystem.components

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.BasicTextField
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.shadow
import androidx.compose.ui.graphics.SolidColor
import androidx.compose.ui.semantics.Role
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.input.VisualTransformation
import androidx.compose.ui.unit.dp
import de.mse.cardhub.designsystem.icons.CardwiseIcon
import de.mse.cardhub.designsystem.icons.CardwiseIconImage
import de.mse.cardhub.designsystem.theme.CardwiseTheme

/** Base input field. Used by [SearchField] and [CardwiseTextField]. */
@Composable
private fun FieldContainer(
    value: String,
    onValueChange: (String) -> Unit,
    placeholder: String,
    modifier: Modifier = Modifier,
    leading: (@Composable () -> Unit)? = null,
    trailing: (@Composable () -> Unit)? = null,
    keyboardOptions: KeyboardOptions = KeyboardOptions.Default,
    visualTransformation: VisualTransformation = VisualTransformation.None,
) {
    val colors = CardwiseTheme.colors
    val shape = CardwiseTheme.shapes.lg
    BasicTextField(
        value = value,
        onValueChange = onValueChange,
        singleLine = true,
        textStyle = CardwiseTheme.typography.body.copy(color = colors.text),
        cursorBrush = SolidColor(colors.blue),
        keyboardOptions = keyboardOptions,
        visualTransformation = visualTransformation,
        modifier = modifier,
        decorationBox = { innerTextField ->
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(46.dp)
                    .shadow(1.dp, shape, ambientColor = ShadowColor, spotColor = ShadowColor)
                    .background(colors.surface, shape)
                    .border(1.dp, colors.line, shape)
                    .padding(start = 13.dp, end = if (trailing != null) 7.dp else 13.dp),
                horizontalArrangement = Arrangement.spacedBy(10.dp),
                verticalAlignment = Alignment.CenterVertically,
            ) {
                leading?.invoke()
                Box(Modifier.weight(1f)) {
                    if (value.isEmpty()) {
                        Text(placeholder, style = CardwiseTheme.typography.body, color = colors.placeholder)
                    }
                    innerTextField()
                }
                trailing?.invoke()
            }
        },
    )
}

/** Search field with search icon and optional filter button. */
@Composable
fun SearchField(
    value: String,
    onValueChange: (String) -> Unit,
    placeholder: String,
    modifier: Modifier = Modifier,
    onFilterClick: (() -> Unit)? = null,
) {
    val colors = CardwiseTheme.colors
    FieldContainer(
        value = value,
        onValueChange = onValueChange,
        placeholder = placeholder,
        modifier = modifier,
        leading = { CardwiseIconImage(CardwiseIcon.Search, tint = colors.muted, size = 18.dp) },
        trailing = onFilterClick?.let { onClick ->
            {
                Box(
                    modifier = Modifier
                        .size(32.dp)
                        .clip(RoundedCornerShape(9.dp))
                        .background(colors.blueSoft)
                        .clickable(role = Role.Button, onClick = onClick),
                    contentAlignment = Alignment.Center,
                ) {
                    CardwiseIconImage(CardwiseIcon.Sliders, tint = colors.blue, size = 18.dp, contentDescription = "Filter öffnen")
                }
            }
        },
    )
}

/**
 * Text field with label, e.g. for login, sign-up and profile.
 * Note: label spacing is not from Figma yet and still needs to be checked.
 */
@Composable
fun CardwiseTextField(
    value: String,
    onValueChange: (String) -> Unit,
    label: String,
    modifier: Modifier = Modifier,
    placeholder: String = "",
    isPassword: Boolean = false,
    keyboardOptions: KeyboardOptions = KeyboardOptions.Default,
) {
    Column(modifier, verticalArrangement = Arrangement.spacedBy(6.dp)) {
        Text(label, style = CardwiseTheme.typography.captionStrong, color = CardwiseTheme.colors.muted)
        FieldContainer(
            value = value,
            onValueChange = onValueChange,
            placeholder = placeholder,
            keyboardOptions = keyboardOptions,
            visualTransformation = if (isPassword) PasswordVisualTransformation() else VisualTransformation.None,
        )
    }
}
