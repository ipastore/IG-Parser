package tree

import (
	"fmt"
	"log"
)

type Statement struct {

	// Regulative Statement
	Attributes                            *Node
	AttributesPropertySimple              *Node
	AttributesPropertyComplex             *Node
	Deontic                               *Node
	Aim                                   *Node
	DirectObject                          *Node
	DirectObjectComplex                   *Node
	DirectObjectPropertySimple            *Node
	DirectObjectPropertyComplex           *Node
	IndirectObject                        *Node
	IndirectObjectComplex                 *Node
	IndirectObjectPropertySimple          *Node
	IndirectObjectPropertyComplex         *Node

	//Constitutive Statement
	ConstitutedEntity                     *Node
	ConstitutedEntityPropertySimple       *Node
	ConstitutedEntityPropertyComplex      *Node
	Modal                                 *Node
	ConstitutiveFunction                  *Node
	ConstitutingProperties                *Node
	ConstitutingPropertiesComplex         *Node
	ConstitutingPropertiesPropertySimple  *Node
	ConstitutingPropertiesPropertyComplex *Node

	// Shared Components
	ActivationConditionSimple             *Node
	ActivationConditionComplex            *Node
	ExecutionConstraintSimple             *Node
	ExecutionConstraintComplex            *Node
	OrElse                                *Node
}

func (s *Statement) String() string {
	out := ""

	out = s.printComponent(out, s.Attributes, ATTRIBUTES, false)
	out = s.printComponent(out, s.AttributesPropertySimple, ATTRIBUTES_PROPERTY, false)
	out = s.printComponent(out, s.AttributesPropertyComplex, ATTRIBUTES_PROPERTY, true)
	out = s.printComponent(out, s.Deontic, DEONTIC, false)
	out = s.printComponent(out, s.Aim, AIM, false)
	out = s.printComponent(out, s.DirectObject, DIRECT_OBJECT, false)
	out = s.printComponent(out, s.DirectObjectComplex, DIRECT_OBJECT, true)
	out = s.printComponent(out, s.DirectObjectPropertySimple, DIRECT_OBJECT_PROPERTY, false)
	out = s.printComponent(out, s.DirectObjectPropertyComplex, DIRECT_OBJECT_PROPERTY, true)
	out = s.printComponent(out, s.IndirectObject, INDIRECT_OBJECT, false)
	out = s.printComponent(out, s.IndirectObjectComplex, INDIRECT_OBJECT, true)
	out = s.printComponent(out, s.IndirectObjectPropertySimple, INDIRECT_OBJECT_PROPERTY, false)
	out = s.printComponent(out, s.IndirectObjectPropertyComplex, INDIRECT_OBJECT_PROPERTY, true)

	out = s.printComponent(out, s.ActivationConditionSimple, ACTIVATION_CONDITION, false)
	out = s.printComponent(out, s.ActivationConditionComplex, ACTIVATION_CONDITION, true)
	out = s.printComponent(out, s.ExecutionConstraintSimple, EXECUTION_CONSTRAINT, false)
	out = s.printComponent(out, s.ExecutionConstraintComplex, EXECUTION_CONSTRAINT, true)

	out = s.printComponent(out, s.ConstitutedEntity, CONSTITUTED_ENTITY, false)
	out = s.printComponent(out, s.ConstitutedEntityPropertySimple, CONSTITUTED_ENTITY_PROPERTY, false)
	out = s.printComponent(out, s.ConstitutedEntityPropertyComplex, CONSTITUTED_ENTITY_PROPERTY, true)
	out = s.printComponent(out, s.Modal, MODAL, false)
	out = s.printComponent(out, s.ConstitutiveFunction, CONSTITUTIVE_FUNCTION, false)
	out = s.printComponent(out, s.ConstitutingProperties, CONSTITUTING_PROPERTIES, false)
	out = s.printComponent(out, s.ConstitutingPropertiesComplex, CONSTITUTING_PROPERTIES, true)
	out = s.printComponent(out, s.ConstitutingPropertiesPropertySimple, CONSTITUTING_PROPERTIES_PROPERTY, false)
	out = s.printComponent(out, s.ConstitutingPropertiesPropertyComplex, CONSTITUTING_PROPERTIES_PROPERTY, true)

	out = s.printComponent(out, s.OrElse, OR_ELSE, true)

	return out
}

/*
Appends component information for output string
Input:
- input string to append output to
- Node whose content is to be appended
- Symbol associated with component
- Indicator whether component is complex

Returns string for printing
*/
func (s *Statement) printComponent(inputString string, node *Node, nodeSymbol string, complex bool) string {

	sep := ": "
	suffix := "\n"
	complexPrefix := "{\n"
	complexSuffix := "\n}"

	// If node is not nil
	if node != nil {
		// Print symbol
		inputString += nodeSymbol + sep
		// Add core content
		if complex {
			// Complex (i.e., nested) node output

			// Append complex node-specific information to the end of nested statement
			// Assumes that suffix and annotations are in string format for nodes that have nested statements
			// TODO: see whether that needs to be adjusted
			if node.Suffix != nil {
				complexSuffix += " (Suffix: " + node.Suffix.(string) + ")"
			}
			if node.Annotations != nil {
				complexSuffix += " (Annotation: " + node.Annotations.(string) + ")"
			}
			if node.PrivateNodeLinks != nil {
				complexSuffix += " (Private links: " + fmt.Sprint(node.PrivateNodeLinks) + ")"
			}
			if node.GetComponentName() != "" {
				complexSuffix += " (Component name: " + fmt.Sprint(node.GetComponentName()) + ")"
			}

			inputString += complexPrefix + node.String() + complexSuffix
		} else {
			// Simple output
			inputString += node.String()
		}
		// Append suffix
		inputString += suffix
	}
	return inputString
}

/*
Stringifies institutional statement
*/
func (s *Statement) Stringify() string {
	log.Fatal("Stringify() is not yet implemented.")
	return ""
}

/*
Generates map of arrays containing pointers to leaf nodes in each component.
Key is an incrementing index, and value is an array of the corresponding nodes.
It further returns an array containing the component keys alongside the number of leaf nodes per component,
in order to reconstruct the linkage between the index in the first return value and the components they relate to.

Example: The first return may include two ATTRIBUTES component trees separated by synthetic AND connections (sAND)
based on different logical combination within the attributes component that are not genuine logical relationships (i.e.,
not signaled using [AND], [OR], or [XOR], but inferred during parsing based on the occurrence of multiple such combinations
within an Attributes component expression (e.g., A((Sellers [AND] Buyers) from (Northern [OR] Southern) states)).
Internally, this would be represented as ((Sellers [AND] Buyers] [sAND] (Northern [OR] Southern))', and returned as separate
trees with index 0 (Sellers [AND] Buyers) and 1 (Northern [OR] Southern).
The second return indicates the fact that the first two entries in the first return type instance are of type ATTRIBUTES by holding
an entry '"ATTRIBUTES": 2', etc.

*/
func (s *Statement) GenerateLeafArrays() ([][]*Node, map[string]int) {
	return s.generateLeafArrays(0)
}

/*
Generates map of arrays containing pointers to leaf nodes in each component.
Key is an incrementing index, and value is an array of the corresponding nodes.
It further returns an array containing the component keys alongside the number of leaf nodes per component,
in order to reconstruct the linkage between the index in the first return value and the components they relate to.

Note: This variant only returns nodes that have a non-nil suffix.

Example: The first return may include two ATTRIBUTES component trees separated by synthetic AND connections (sAND)
based on different logical combination within the attributes component that are not genuine logical relationships (i.e.,
not signaled using [AND], [OR], or [XOR], but inferred during parsing based on the occurrence of multiple such combinations
within an Attributes component expression (e.g., A((Sellers [AND] Buyers) from (Northern [OR] Southern) states)).
Internally, this would be represented as ((Sellers [AND] Buyers] [sAND] (Northern [OR] Southern))', and returned as separate
trees with index 0 (Sellers [AND] Buyers) and 1 (Northern [OR] Southern).
The second return indicates the fact that the first two entries in the first return type instance are of type ATTRIBUTES by holding
an entry '"ATTRIBUTES": 2', etc.

*/
func (s *Statement) GenerateLeafArraysSuffixOnly() ([][]*Node, map[string]int) {
	return s.generateLeafArrays(1)
}

/*
Generates map of arrays containing pointers to leaf nodes in each component.
Key is an incrementing index, and value is an array of the corresponding nodes.
It further returns an array containing the component keys alongside the number of leaf nodes per component,
in order to reconstruct the linkage between the index in the first return value and the components they relate to.

Input: level indicates selection of nodes considered in aggregation (0 --> all nodes, 1 --> nodes with non-nil suffix only)

Example: The first return may include two ATTRIBUTES component trees separated by synthetic AND connections (sAND)
based on different logical combination within the attributes component that are not genuine logical relationships (i.e.,
not signaled using [AND], [OR], or [XOR], but inferred during parsing based on the occurrence of multiple such combinations
within an Attributes component expression (e.g., A((Sellers [AND] Buyers) from (Northern [OR] Southern) states)).
Internally, this would be represented as ((Sellers [AND] Buyers] [sAND] (Northern [OR] Southern))', and returned as separate
trees with index 0 (Sellers [AND] Buyers) and 1 (Northern [OR] Southern).
The second return indicates the fact that the first two entries in the first return type instance are of type ATTRIBUTES by holding
an entry '"ATTRIBUTES": 2', etc.

*/
func (s *Statement) generateLeafArrays(level int) ([][]*Node, map[string]int) {

	// Map holding reference from component type (e.g., ATTRIBUTES) to number of entries (relevant for reconstruction)
	referenceMap := map[string]int{}

	// Counter for overall number of entries
	nodesMap := make([][]*Node, 0)

	// Regulative components
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.Attributes, ATTRIBUTES, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.AttributesPropertySimple, ATTRIBUTES_PROPERTY, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.AttributesPropertyComplex, ATTRIBUTES_PROPERTY_REFERENCE, true, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.Deontic, DEONTIC, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.Aim, AIM, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.DirectObject, DIRECT_OBJECT, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.DirectObjectComplex, DIRECT_OBJECT_REFERENCE, true, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.DirectObjectPropertySimple, DIRECT_OBJECT_PROPERTY, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.DirectObjectPropertyComplex, DIRECT_OBJECT_PROPERTY_REFERENCE, true, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.IndirectObject, INDIRECT_OBJECT, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.IndirectObjectComplex, INDIRECT_OBJECT_REFERENCE, true, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.IndirectObjectPropertySimple, INDIRECT_OBJECT_PROPERTY, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.IndirectObjectPropertyComplex, INDIRECT_OBJECT_PROPERTY_REFERENCE, true, level)

	// Context
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ActivationConditionSimple, ACTIVATION_CONDITION, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ActivationConditionComplex, ACTIVATION_CONDITION_REFERENCE, true, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ExecutionConstraintSimple, EXECUTION_CONSTRAINT, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ExecutionConstraintComplex, EXECUTION_CONSTRAINT_REFERENCE, true, level)

	// Constitutive components
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutedEntity, CONSTITUTED_ENTITY, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutedEntityPropertySimple, CONSTITUTED_ENTITY_PROPERTY, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutedEntityPropertyComplex, CONSTITUTED_ENTITY_PROPERTY_REFERENCE, true, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.Modal, MODAL, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutiveFunction, CONSTITUTIVE_FUNCTION, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutingProperties, CONSTITUTING_PROPERTIES, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutingPropertiesComplex, CONSTITUTING_PROPERTIES_REFERENCE, true, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutingPropertiesPropertySimple, CONSTITUTING_PROPERTIES_PROPERTY, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutingPropertiesPropertyComplex, CONSTITUTING_PROPERTIES_PROPERTY_REFERENCE, true, level)

	// Shared components
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.OrElse, OR_ELSE, true, level)

	return nodesMap, referenceMap
}

/*
Generates a leaf array for a given component under consideration of node as being of simple or complex nature.
Appends to existing structure if provided (i.e., not nil) to allow for iterative invocation.
For a version that allows for iterative invocation, consider #getComponentLeafArray.
For returning only leaves that contain suffix information, consider #getComponentLeafArrayWithSuffix.

Input:
- Reference to component node for which leaf elements are to be extracted
- Component symbol associated with component
- Indicator whether element embedded in node is complex (i.e., nested statement)
- Indicator whether all leaf nodes should be returned, or only one satisfying particular conditions
  (0 --> all nodes, 1 --> only ones with non-empty suffix).

Returns:
- Node map of nodes associated with components
- Reference map counting number of components
*/
func GetSingleComponentLeafArray(componentNode *Node, componentSymbol string, complex bool, level int) ([][]*Node, map[string]int) {

	// Map holding reference from component type (e.g., ATTRIBUTES) to number of entries (relevant for reconstruction)
	referenceMap := map[string]int{}

	// Counter for overall number of entries
	nodesMap := make([][]*Node, 0)

	return getComponentLeafArray(nodesMap, referenceMap, componentNode, componentSymbol, complex, level)
}

/*
Generates a leaf array for a given component under consideration of node as being of simple or complex nature.
Appends to existing structure if provided (i.e., not nil) to allow for iterative invocation.
For returning only leaves that contain suffix information consider #getComponentLeafArrayWithSuffix.
Input:
- maps of nodes potentially including existing nodes for other components. Will be created internally if nil
  (to allow iterative invocation).
- reference map that indexes the number of nodes associated with a specific component (to retain association).
  Will be created internally if nil (to allow iterative invocation).
- Reference to component node for which leaf elements are to be extracted
- Component symbol associated with component
- Indicator whether element embedded in node is complex (i.e., nested statement)
- Indicator whether all leaf nodes should be returned, or only one satisfying particular conditions
  (0 --> all nodes, 1 --> only ones with non-empty suffix).

Returns:
- Node map of nodes associated with components
- Reference map counting number of components
*/
func getComponentLeafArray(nodesMap [][]*Node, referenceMap map[string]int, componentNode *Node, componentSymbol string, complex bool, level int) ([][]*Node, map[string]int) {

	if componentNode == nil {
		fmt.Println("No component node found - returning unmodified node and reference map ...")
		return nodesMap, referenceMap
	}

	// Initialize data structures if nil
	if nodesMap == nil {
		nodesMap = make([][]*Node, 1)
	}

	if referenceMap == nil {
		referenceMap = make(map[string]int, 1)
	}

	// Check for complex content
	if complex {
		// Embed nested statement in node structure, before adding to node map
		nodesMap = append(nodesMap, []*Node{componentNode})

		// since statements can be combined, they are returned as a single element
		referenceMap[componentSymbol] = 1
	} else {
		// Counter for number of elements in given component
		i := 0
		// Add array of leaf nodes attached to general array
		for _, v := range componentNode.GetLeafNodes() {
			nodesMap = append(nodesMap, v)
			i++
		}
		// Add number of nodes referring to a particular component
		referenceMap[componentSymbol] = i
	}

	// Return modified or generated structures
	return nodesMap, referenceMap
}
